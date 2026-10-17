package mongo

import (
	"context"
	"errors"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// userRepo is the MongoDB implementation of repository.UserRepository.
type userRepo struct {
	collection *mongo.Collection
}

// NewUserRepository creates a UserRepository backed by the "users" collection.
// EnsureIndexes must have run so that emails are unique.
func NewUserRepository(db *mongo.Database) repository.UserRepository {
	return &userRepo{collection: db.Collection(usersCollection)}
}

// Create inserts a new user.
func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return repository.ErrDuplicate
	}
	return err
}

// GetByID finds a user by ID.
func (r *userRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return r.findOne(ctx, bson.M{"_id": id})
}

// GetByEmail finds a user by email.
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.findOne(ctx, bson.M{"email": email})
}

// findOne decodes the first user matching filter.
func (r *userRepo) findOne(ctx context.Context, filter bson.M) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
