package service

import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/realtime"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Field limits for tasks.
const (
	maxTaskTitle       = 200
	maxTaskDescription = 2000
)

// TaskInput is what a client sends to create or update a task.
type TaskInput struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Priority    models.Priority      `json:"priority"`
	Completed   bool                 `json:"completed"`
	DeadlineIDs []primitive.ObjectID `json:"deadline_ids"`
	TagIDs      []primitive.ObjectID `json:"tag_ids"`
	// Version is the version the client last saw. Required for updates.
	Version int64 `json:"version"`
}

// TaskService manages a user's tasks.
type TaskService struct {
	tasks     repository.TaskRepository
	deadlines repository.DeadlineRepository
	tags      repository.TagRepository
	publisher Publisher
	now       clock
}

// NewTaskService creates a TaskService.
func NewTaskService(tasks repository.TaskRepository, deadlines repository.DeadlineRepository, tags repository.TagRepository, publisher Publisher) *TaskService {
	return &TaskService{tasks: tasks, deadlines: deadlines, tags: tags, publisher: publisher, now: systemClock}
}

// List returns all of the user's tasks.
func (s *TaskService) List(ctx context.Context, userID primitive.ObjectID) ([]*models.Task, error) {
	return s.tasks.List(ctx, userID)
}

// Get returns one of the user's tasks.
func (s *TaskService) Get(ctx context.Context, userID, id primitive.ObjectID) (*models.Task, error) {
	return s.tasks.Get(ctx, userID, id)
}

// Create validates and stores a new task.
func (s *TaskService) Create(ctx context.Context, userID primitive.ObjectID, in TaskInput) (*models.Task, error) {
	task := &models.Task{Meta: newMeta(userID, s.now())}
	if err := s.apply(ctx, userID, task, in); err != nil {
		return nil, err
	}
	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, err
	}
	s.publish(userID, realtime.TaskCreated, task)
	return task, nil
}

// Update validates and replaces a task. It returns ErrConflict when the
// task changed since in.Version.
func (s *TaskService) Update(ctx context.Context, userID, id primitive.ObjectID, in TaskInput) (*models.Task, error) {
	if err := checkVersion(in.Version); err != nil {
		return nil, err
	}
	task, err := s.tasks.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if err := s.apply(ctx, userID, task, in); err != nil {
		return nil, err
	}
	task.Version = in.Version
	task.UpdatedAt = s.now()
	if err := s.tasks.Update(ctx, task); err != nil {
		return nil, err
	}
	s.publish(userID, realtime.TaskUpdated, task)
	return task, nil
}

// Delete removes a task.
func (s *TaskService) Delete(ctx context.Context, userID, id primitive.ObjectID) error {
	if err := s.tasks.Delete(ctx, userID, id); err != nil {
		return err
	}
	s.publisher.Publish(userID, realtime.Event{Type: realtime.TaskDeleted, ID: id.Hex()})
	return nil
}

// apply validates in and copies it onto task.
func (s *TaskService) apply(ctx context.Context, userID primitive.ObjectID, task *models.Task, in TaskInput) error {
	title, err := cleanText("title", in.Title, true, maxTaskTitle)
	if err != nil {
		return err
	}
	description, err := cleanText("description", in.Description, false, maxTaskDescription)
	if err != nil {
		return err
	}
	if !in.Priority.Valid() {
		return invalid("priority", "must be low, medium or high")
	}
	deadlineIDs, err := cleanIDs(ctx, s.deadlines, userID, "deadline_ids", in.DeadlineIDs)
	if err != nil {
		return err
	}
	tagIDs, err := cleanIDs(ctx, s.tags, userID, "tag_ids", in.TagIDs)
	if err != nil {
		return err
	}

	task.Title = title
	task.Description = description
	task.Priority = in.Priority
	task.Completed = in.Completed
	task.DeadlineIDs = deadlineIDs
	task.TagIDs = tagIDs
	return nil
}

// publish sends a created/updated event for task.
func (s *TaskService) publish(userID primitive.ObjectID, eventType string, task *models.Task) {
	s.publisher.Publish(userID, realtime.Event{Type: eventType, ID: task.ID.Hex(), Data: task})
}

// publishChanged reloads the given tasks and tells devices they changed.
// It is used after a deadline or tag was unlinked from them.
func (s *TaskService) publishChanged(ctx context.Context, userID primitive.ObjectID, ids []primitive.ObjectID) {
	for _, id := range ids {
		if task, err := s.tasks.Get(ctx, userID, id); err == nil {
			s.publish(userID, realtime.TaskUpdated, task)
		}
	}
}
