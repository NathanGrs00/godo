package _interface

import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	// Create adds a new user to the database
	Create(ctx context.Context, user *models.User) error
	// GetByID retrieves a user by its ID
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	// GetByEmail retrieves a user by its email
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	// Update modifies an existing user
	Update(ctx context.Context, user *models.User) error
	// Delete removes a user by its ID
	Delete(ctx context.Context, id primitive.ObjectID) error
}
