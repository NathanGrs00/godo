package mongo

import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	_interface "github.com/NathanGrs00/godo/internal/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// userRepo represents the MongoDB implementation of the User repository.
type userRepo struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new instance of UserRepository using the provided MongoDB database.
func NewUserRepository(db *mongo.Database) _interface.UserRepository {
	return &userRepo{collection: db.Collection("users")}
}

// Create adds a new user to the MongoDB collection.
// r is like self in Python, it refers to the instance of userRepo.
func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// GetById retrieves a user by its ID from the MongoDB collection.
func (r *userRepo) GetById(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	// user will hold the result of the query
	var user models.User
	// Filter to find the user by its ID
	filter := bson.M{"_id": id}
	// Perform the query and decode the result into the user variable
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	// Return nil for the user and the error if there was an issue
	if err != nil {
		return nil, err
	}
	// Return the found user
	return &user, nil
}

// GetByEmail finds a user by the given email.
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	// Make a user struct
	var user models.User
	// Make a filter by email
	filter := bson.M{"email": email}
	// Perform the query and decode it into a user variable
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	// If any errors occur, return them
	if err != nil {
		return nil, err
	}
	// Return the user
	return &user, nil
}

// Update modifies an existing user in the MongoDB collection.
func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	// Filter to find the user by its ID
	filter := bson.M{"_id": user.ID}
	// Update the user document with the new values
	update := bson.D{{"$set", user}}
	// Perform the update operation
	_, err := r.collection.UpdateOne(ctx, filter, update)
	// Return the error if any
	return err
}

// Delete removes an existing user in the MongoDB collection.
func (r *userRepo) Delete(ctx context.Context, userID primitive.ObjectID) error {
	// Filter to find the user by ID
	filter := bson.M{"_id": userID}
	// Perform the update operation
	_, err := r.collection.DeleteOne(ctx, filter)
	// Return the error if any
	return err
}
