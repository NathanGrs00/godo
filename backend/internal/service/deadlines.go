package service

import (
	"context"
	"time"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/realtime"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Field limits for deadlines.
const (
	maxDeadlineTitle       = 200
	maxDeadlineDescription = 2000
)

// DeadlineInput is what a client sends to create or update a deadline.
type DeadlineInput struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Date        time.Time            `json:"date"`
	TagIDs      []primitive.ObjectID `json:"tag_ids"`
	// Version is the version the client last saw. Required for updates.
	Version int64 `json:"version"`
}

// DeadlineService manages a user's deadlines.
type DeadlineService struct {
	deadlines repository.DeadlineRepository
	tags      repository.TagRepository
	tasks     *TaskService
	publisher Publisher
	now       clock
}

// NewDeadlineService creates a DeadlineService. tasks is used to unlink and
// re-announce tasks when a deadline is deleted.
func NewDeadlineService(deadlines repository.DeadlineRepository, tags repository.TagRepository, tasks *TaskService, publisher Publisher) *DeadlineService {
	return &DeadlineService{deadlines: deadlines, tags: tags, tasks: tasks, publisher: publisher, now: systemClock}
}

// List returns all of the user's deadlines.
func (s *DeadlineService) List(ctx context.Context, userID primitive.ObjectID) ([]*models.Deadline, error) {
	deadlines, err := s.deadlines.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	now := s.now()
	for _, d := range deadlines {
		markPassed(d, now)
	}
	return deadlines, nil
}

// Get returns one of the user's deadlines.
func (s *DeadlineService) Get(ctx context.Context, userID, id primitive.ObjectID) (*models.Deadline, error) {
	deadline, err := s.deadlines.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	markPassed(deadline, s.now())
	return deadline, nil
}

// Create validates and stores a new deadline.
func (s *DeadlineService) Create(ctx context.Context, userID primitive.ObjectID, in DeadlineInput) (*models.Deadline, error) {
	deadline := &models.Deadline{Meta: newMeta(userID, s.now())}
	if err := s.apply(ctx, userID, deadline, in); err != nil {
		return nil, err
	}
	if err := s.deadlines.Create(ctx, deadline); err != nil {
		return nil, err
	}
	s.publish(userID, realtime.DeadlineCreated, deadline)
	return deadline, nil
}

// Update validates and replaces a deadline. It returns ErrConflict when the
// deadline changed since in.Version.
func (s *DeadlineService) Update(ctx context.Context, userID, id primitive.ObjectID, in DeadlineInput) (*models.Deadline, error) {
	if err := checkVersion(in.Version); err != nil {
		return nil, err
	}
	deadline, err := s.deadlines.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if err := s.apply(ctx, userID, deadline, in); err != nil {
		return nil, err
	}
	deadline.Version = in.Version
	deadline.UpdatedAt = s.now()
	if err := s.deadlines.Update(ctx, deadline); err != nil {
		return nil, err
	}
	s.publish(userID, realtime.DeadlineUpdated, deadline)
	return deadline, nil
}

// Delete removes a deadline and unlinks it from every task.
func (s *DeadlineService) Delete(ctx context.Context, userID, id primitive.ObjectID) error {
	if err := s.deadlines.Delete(ctx, userID, id); err != nil {
		return err
	}
	s.publisher.Publish(userID, realtime.Event{Type: realtime.DeadlineDeleted, ID: id.Hex()})

	changed, err := s.tasks.tasks.RemoveDeadline(ctx, userID, id)
	if err != nil {
		return err
	}
	s.tasks.publishChanged(ctx, userID, changed)
	return nil
}

// apply validates in and copies it onto deadline.
func (s *DeadlineService) apply(ctx context.Context, userID primitive.ObjectID, deadline *models.Deadline, in DeadlineInput) error {
	title, err := cleanText("title", in.Title, true, maxDeadlineTitle)
	if err != nil {
		return err
	}
	description, err := cleanText("description", in.Description, false, maxDeadlineDescription)
	if err != nil {
		return err
	}
	if in.Date.IsZero() {
		return invalid("date", "is required")
	}
	tagIDs, err := cleanIDs(ctx, s.tags, userID, "tag_ids", in.TagIDs)
	if err != nil {
		return err
	}

	deadline.Title = title
	deadline.Description = description
	deadline.Date = in.Date.UTC().Truncate(time.Millisecond)
	deadline.TagIDs = tagIDs
	markPassed(deadline, s.now())
	return nil
}

// publish sends a created/updated event for deadline.
func (s *DeadlineService) publish(userID primitive.ObjectID, eventType string, deadline *models.Deadline) {
	s.publisher.Publish(userID, realtime.Event{Type: eventType, ID: deadline.ID.Hex(), Data: deadline})
}

// publishChanged reloads the given deadlines and tells devices they changed.
func (s *DeadlineService) publishChanged(ctx context.Context, userID primitive.ObjectID, ids []primitive.ObjectID) {
	for _, id := range ids {
		if deadline, err := s.Get(ctx, userID, id); err == nil {
			s.publish(userID, realtime.DeadlineUpdated, deadline)
		}
	}
}

// markPassed sets Passed when the deadline's date is before now.
func markPassed(deadline *models.Deadline, now time.Time) {
	deadline.Passed = deadline.Date.Before(now)
}
