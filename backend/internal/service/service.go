// Package service holds the application's business rules: validation,
// ownership, linking tasks to deadlines and tags, and telling a user's other
// devices about every change.
//
// HTTP handlers call services; services call repositories.
package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/realtime"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Errors returned by services. Handlers turn them into HTTP status codes.
var (
	// ErrNotFound means the item does not exist or belongs to someone else.
	ErrNotFound = repository.ErrNotFound
	// ErrConflict means the item was changed (for example on another device)
	// since the client loaded it. The client should reload and retry.
	ErrConflict = repository.ErrConflict
	// ErrEmailTaken means an account with that email already exists.
	ErrEmailTaken = errors.New("email already registered")
	// ErrInvalidCredentials means the email or password is wrong.
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// ValidationError means the input was rejected. Field names the bad field.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s %s", e.Field, e.Message)
}

// invalid is a shorthand for creating a ValidationError.
func invalid(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}

// Publisher sends change events to a user's connected devices.
// realtime.Hub implements it.
type Publisher interface {
	Publish(userID primitive.ObjectID, event realtime.Event)
}

// clock returns the current time. Services keep one so tests can fix time.
type clock func() time.Time

// systemClock returns the current UTC time, rounded to milliseconds because
// that is the precision MongoDB stores.
func systemClock() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}

// newMeta fills in the shared fields of a new document.
func newMeta(userID primitive.ObjectID, now time.Time) models.Meta {
	return models.Meta{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// cleanText trims a text field and checks its length.
// When required is true an empty value is rejected.
func cleanText(field, value string, required bool, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if required && value == "" {
		return "", invalid(field, "is required")
	}
	if len([]rune(value)) > maxLen {
		return "", invalid(field, fmt.Sprintf("must be at most %d characters", maxLen))
	}
	return value, nil
}

// cleanIDs removes duplicate IDs and checks that every ID refers to one of
// the user's documents in store. It never returns nil, so JSON shows [].
func cleanIDs[T models.Entity](ctx context.Context, store repository.Store[T], userID primitive.ObjectID, field string, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
	unique := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if !slices.Contains(unique, id) {
			unique = append(unique, id)
		}
	}
	ok, err := store.AllExist(ctx, userID, unique)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, invalid(field, "contains an unknown id")
	}
	return unique, nil
}

// checkVersion rejects updates that do not say which version they change.
func checkVersion(version int64) error {
	if version < 1 {
		return invalid("version", "is required when updating")
	}
	return nil
}
