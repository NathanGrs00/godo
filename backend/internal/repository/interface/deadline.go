package _interface

// Import necessary packages
import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeadlineRepository defines the interface for deadline-related database operations
type DeadlineRepository interface {
	// Create adds a new deadline to the database
	Create(ctx context.Context, deadline *models.Deadline) error
	// GetById retrieves a deadline by its ID
	GetById(ctx context.Context, id primitive.ObjectID) (*models.Deadline, error)
	// GetByUserId retrieves all deadlines for a specific user
	GetByUserId(ctx context.Context, userId primitive.ObjectID) ([]*models.Deadline, error)
	// Update modifies an existing deadline
	Update(ctx context.Context, deadline *models.Deadline) error
	// Delete removes a deadline by its ID
	Delete(ctx context.Context, id primitive.ObjectID) error
	// GetAll retrieves all deadlines
	GetAll(ctx context.Context) ([]*models.Deadline, error)
}
