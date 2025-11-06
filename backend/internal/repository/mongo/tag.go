package mongo

import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	_interface "github.com/NathanGrs00/godo/internal/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// tagRepo represents the MongoDB implementation of the Tag repository.
type tagRepo struct {
	collection *mongo.Collection
}

// NewTagRepository creates a new instance of TagRepository using the provided MongoDB database.
func NewTagRepository(db *mongo.Database) _interface.TagRepository {
	return &tagRepo{collection: db.Collection("tags")}
}

// Create adds a new tag to the MongoDB collection.
// r is like self in Python, it refers to the instance of tagRepo.
func (r *tagRepo) Create(ctx context.Context, tag *models.Tag) error {
	_, err := r.collection.InsertOne(ctx, tag)
	return err
}

// GetById retrieves a tag by its ID from the MongoDB collection.
func (r *tagRepo) GetById(ctx context.Context, id primitive.ObjectID) (*models.Tag, error) {
	// tag will hold the result of the query
	var tag models.Tag
	// Filter to find the tag by its ID
	filter := bson.M{"_id": id}
	// Perform the query and decode the result into the tag variable
	err := r.collection.FindOne(ctx, filter).Decode(&tag)
	// Return nil for the tag and the error if there was an issue
	if err != nil {
		return nil, err
	}
	// Return the found tag
	return &tag, nil
}

// GetByUserId retrieves all tags for a specific user from the MongoDB collection.
func (r *tagRepo) GetByUserId(ctx context.Context, userID primitive.ObjectID) ([]*models.Tag, error) {
	// Filter to find tags by user ID
	filter := bson.M{"user_id": userID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed after processing
	defer cursor.Close(ctx)

	// tags will hold the list of tags for the user
	var tags []*models.Tag
	// Iterate through the cursor
	for cursor.Next(ctx) {
		// Decode each document into a Tag struct
		var tag models.Tag
		// If error occurs during decoding, return it
		if err := cursor.Decode(&tag); err != nil {
			return nil, err
		}
		// Append the tag to the tags slice
		tags = append(tags, &tag)
	}

	// Return the list of tags
	return tags, nil
}

// Update modifies an existing tag in the MongoDB collection.
func (r *tagRepo) Update(ctx context.Context, tag *models.Tag) error {
	// Filter to find the tag by its ID
	filter := bson.M{"_id": tag.ID}
	// Update the tag document with the new values
	update := bson.D{{"$set", tag}}
	// Perform the update operation
	_, err := r.collection.UpdateOne(ctx, filter, update)
	// Return the error if any
	return err
}

// Delete removes an existing tag in the MongoDB collection.
func (r *tagRepo) Delete(ctx context.Context, tagID primitive.ObjectID) error {
	// Filter to find the tag by ID
	filter := bson.M{"_id": tagID}
	// Perform the update operation
	_, err := r.collection.DeleteOne(ctx, filter)
	// Return the error if any
	return err
}

func (r *tagRepo) GetAll(ctx context.Context) ([]*models.Tag, error) {
	// tags will hold the list of all tags
	var tags []*models.Tag

	// Find all documents in the collection
	cursor, err := r.collection.Find(ctx, bson.M{})
	// If error occurs during Find, return it
	if err != nil {
		// return empty slice and the error
		return nil, err
	}
	// Ensure the cursor is closed after processing
	defer cursor.Close(ctx)

	// Iterate through the cursor
	for cursor.Next(ctx) {
		// Decode each document into a Tag struct
		var tag models.Tag
		// If error occurs during decoding, return it
		if err := cursor.Decode(&tag); err != nil {
			return nil, err
		}
		// Append the tag to the tags slice
		tags = append(tags, &tag)
	}
	// Return the list of tags
	return tags, nil
}
