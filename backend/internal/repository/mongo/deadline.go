package mongo

import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	_interface "github.com/NathanGrs00/godo/internal/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// deadlineRepo represents the MongoDB implementation of the Deadline repository.
type deadlineRepo struct {
	collection *mongo.Collection
}

// NewDeadlineRepository creates a new instance of DeadlineRepository using the provided MongoDB database.
func NewDeadlineRepository(db *mongo.Database) _interface.DeadlineRepository {
	return &deadlineRepo{collection: db.Collection("deadlines")}
}

// Create adds a new deadline to the MongoDB collection.
// r is like self in Python, it refers to the instance of deadlineRepo.
func (r *deadlineRepo) Create(ctx context.Context, deadline *models.Deadline) error {
	_, err := r.collection.InsertOne(ctx, deadline)
	return err
}

// GetById retrieves a deadline by its ID from the MongoDB collection.
func (r *deadlineRepo) GetById(ctx context.Context, id primitive.ObjectID) (*models.Deadline, error) {
	// deadline will hold the result of the query
	var deadline models.Deadline
	// Filter to find the deadline by its ID
	filter := bson.M{"_id": id}
	// Perform the query and decode the result into the deadline variable
	err := r.collection.FindOne(ctx, filter).Decode(&deadline)
	// Return nil for the deadline and the error if there was an issue
	if err != nil {
		return nil, err
	}
	// Return the found deadline
	return &deadline, nil
}

// GetByUserId retrieves all deadlines for a specific user from the MongoDB collection.
func (r *deadlineRepo) GetByUserId(ctx context.Context, userID primitive.ObjectID) ([]*models.Deadline, error) {
	// Filter to find deadlines by user ID
	filter := bson.M{"user_id": userID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed after processing
	defer cursor.Close(ctx)

	// deadlines will hold the list of deadlines for the user
	var deadlines []*models.Deadline
	// Iterate through the cursor
	for cursor.Next(ctx) {
		// Decode each document into a Deadline struct
		var deadline models.Deadline
		// If error occurs during decoding, return it
		if err := cursor.Decode(&deadline); err != nil {
			return nil, err
		}
		// Append the deadline to the deadlines slice
		deadlines = append(deadlines, &deadline)
	}

	// Return the list of deadlines
	return deadlines, nil
}

// Update modifies an existing deadline in the MongoDB collection.
func (r *deadlineRepo) Update(ctx context.Context, deadline *models.Deadline) error {
	// Filter to find the deadline by its ID
	filter := bson.M{"_id": deadline.ID}
	// Update the deadline document with the new values
	update := bson.D{{"$set", deadline}}
	// Perform the update operation
	_, err := r.collection.UpdateOne(ctx, filter, update)
	// Return the error if any
	return err
}

// Delete removes an existing deadline in the MongoDB collection.
func (r *deadlineRepo) Delete(ctx context.Context, deadlineID primitive.ObjectID) error {
	// Filter to find the deadline by ID
	filter := bson.M{"_id": deadlineID}
	// Perform the update operation
	_, err := r.collection.DeleteOne(ctx, filter)
	// Return the error if any
	return err
}
