package _interface

// Import necessary packages
import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TagRepository defines the interface for tag-related database operations
type TagRepository interface {
	// Create adds a new tag to the database
	Create(ctx context.Context, tag *models.Tag) error
	// GetByID retrieves a tag by its ID
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Tag, error)
	// GetByUserID retrieves all tags for a specific user
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Tag, error)
	// Update modifies an existing tag
	Update(ctx context.Context, tag *models.Tag) error
	// Delete removes a tag by its ID
	Delete(ctx context.Context, id primitive.ObjectID) error
}
