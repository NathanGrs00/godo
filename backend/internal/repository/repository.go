// Package repository defines how the application stores its data.
//
// Services depend only on these interfaces. The mongo sub-package is the
// production implementation; the memory sub-package is an in-memory
// implementation used by tests.
package repository

import (
	"context"
	"errors"

	"github.com/NathanGrs00/godo/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Errors returned by every repository implementation.
var (
	// ErrNotFound means no document with that ID exists for that user.
	ErrNotFound = errors.New("not found")
	// ErrConflict means the document was changed since the caller read it
	// (its version no longer matches).
	ErrConflict = errors.New("version conflict")
	// ErrDuplicate means a unique field (such as a user's email) is taken.
	ErrDuplicate = errors.New("duplicate")
)

// Store holds documents that belong to a user. Every read and write is
// scoped to a user ID, so one user can never see or change another's data.
type Store[T models.Entity] interface {
	// Create stores a new document. Its Meta must already be filled in.
	Create(ctx context.Context, item T) error
	// Get returns one document, or ErrNotFound.
	Get(ctx context.Context, userID, id primitive.ObjectID) (T, error)
	// List returns all of a user's documents, oldest first.
	List(ctx context.Context, userID primitive.ObjectID) ([]T, error)
	// Update replaces a document. item's Version must equal the stored
	// version, otherwise ErrConflict is returned. On success the stored
	// version (and item's Version) is increased by one.
	Update(ctx context.Context, item T) error
	// Delete removes a document, or returns ErrNotFound.
	Delete(ctx context.Context, userID, id primitive.ObjectID) error
	// AllExist reports whether every ID belongs to an existing document of the user.
	AllExist(ctx context.Context, userID primitive.ObjectID, ids []primitive.ObjectID) (bool, error)
}

// TaskRepository stores tasks.
type TaskRepository interface {
	Store[*models.Task]
	// RemoveTag unlinks a tag from all of a user's tasks and returns the IDs of the tasks that changed.
	RemoveTag(ctx context.Context, userID, tagID primitive.ObjectID) ([]primitive.ObjectID, error)
	// RemoveDeadline unlinks a deadline from all of a user's tasks and returns the IDs of the tasks that changed.
	RemoveDeadline(ctx context.Context, userID, deadlineID primitive.ObjectID) ([]primitive.ObjectID, error)
}

// DeadlineRepository stores deadlines.
type DeadlineRepository interface {
	Store[*models.Deadline]
	// RemoveTag unlinks a tag from all of a user's deadlines and returns the IDs of the deadlines that changed.
	RemoveTag(ctx context.Context, userID, tagID primitive.ObjectID) ([]primitive.ObjectID, error)
}

// TagRepository stores tags.
type TagRepository interface {
	Store[*models.Tag]
}

// UserRepository stores user accounts.
type UserRepository interface {
	// Create stores a new user, or returns ErrDuplicate if the email is taken.
	Create(ctx context.Context, user *models.User) error
	// GetByID returns a user, or ErrNotFound.
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	// GetByEmail returns a user, or ErrNotFound.
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
