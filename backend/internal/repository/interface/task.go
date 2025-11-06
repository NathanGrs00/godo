package _interface

// Import necessary packages
import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskRepository defines the interface for task-related database operations
type TaskRepository interface {
	// Create adds a new task to the database
	Create(ctx context.Context, task *models.Task) error
	// GetById retrieves a task by its ID
	GetById(ctx context.Context, id primitive.ObjectID) (*models.Task, error)
	// GetByUserId retrieves all tasks for a specific user
	GetByUserId(ctx context.Context, userID primitive.ObjectID) ([]*models.Task, error)
	// Update modifies an existing task
	Update(ctx context.Context, task *models.Task) error
	// Delete removes a task by its ID
	Delete(ctx context.Context, id primitive.ObjectID) error
}
