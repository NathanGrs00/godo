package mongo

import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// taskRepo is the MongoDB implementation of repository.TaskRepository.
type taskRepo struct {
	store[models.Task, *models.Task]
}

// NewTaskRepository creates a TaskRepository backed by the "tasks" collection.
func NewTaskRepository(db *mongo.Database) repository.TaskRepository {
	return &taskRepo{store[models.Task, *models.Task]{collection: db.Collection(tasksCollection)}}
}

// RemoveTag unlinks a tag from every task of the user.
func (r *taskRepo) RemoveTag(ctx context.Context, userID, tagID primitive.ObjectID) ([]primitive.ObjectID, error) {
	return r.pullReference(ctx, userID, "tag_ids", tagID)
}

// RemoveDeadline unlinks a deadline from every task of the user.
func (r *taskRepo) RemoveDeadline(ctx context.Context, userID, deadlineID primitive.ObjectID) ([]primitive.ObjectID, error) {
	return r.pullReference(ctx, userID, "deadline_ids", deadlineID)
}
