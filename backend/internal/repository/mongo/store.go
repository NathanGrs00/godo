// Package mongo implements the repository interfaces on MongoDB.
package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// store is a generic, user-scoped collection. T is the model struct and
// PT is its pointer type, which implements models.Entity.
type store[T any, PT interface {
	*T
	models.Entity
}] struct {
	collection *mongo.Collection
}

// Create inserts a new document.
func (s *store[T, PT]) Create(ctx context.Context, item PT) error {
	_, err := s.collection.InsertOne(ctx, item)
	return err
}

// Get finds one document by ID that belongs to the user.
func (s *store[T, PT]) Get(ctx context.Context, userID, id primitive.ObjectID) (PT, error) {
	var item T
	err := s.collection.FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(&item)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// List returns all of the user's documents, oldest first.
func (s *store[T, PT]) List(ctx context.Context, userID primitive.ObjectID) ([]PT, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed after processing
	defer cursor.Close(ctx)

	// Start with an empty (not nil) slice so the API returns [] instead of null
	items := []PT{}
	for cursor.Next(ctx) {
		var item T
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, cursor.Err()
}

// Update replaces the document if its stored version still matches.
func (s *store[T, PT]) Update(ctx context.Context, item PT) error {
	meta := item.GetMeta()
	expected := meta.Version
	meta.Version = expected + 1

	// Only match the document if nobody changed it since the caller read it
	filter := bson.M{"_id": meta.ID, "user_id": meta.UserID, "version": expected}
	res, err := s.collection.ReplaceOne(ctx, filter, item)
	if err != nil {
		meta.Version = expected
		return err
	}
	if res.MatchedCount == 1 {
		return nil
	}

	// Nothing matched: either the document is gone or it has a newer version
	meta.Version = expected
	count, err := s.collection.CountDocuments(ctx, bson.M{"_id": meta.ID, "user_id": meta.UserID})
	if err != nil {
		return err
	}
	if count == 0 {
		return repository.ErrNotFound
	}
	return repository.ErrConflict
}

// Delete removes the document if it belongs to the user.
func (s *store[T, PT]) Delete(ctx context.Context, userID, id primitive.ObjectID) error {
	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// AllExist checks that every ID refers to one of the user's documents.
func (s *store[T, PT]) AllExist(ctx context.Context, userID primitive.ObjectID, ids []primitive.ObjectID) (bool, error) {
	unique := uniqueIDs(ids)
	if len(unique) == 0 {
		return true, nil
	}
	count, err := s.collection.CountDocuments(ctx, bson.M{"_id": bson.M{"$in": unique}, "user_id": userID})
	if err != nil {
		return false, err
	}
	return count == int64(len(unique)), nil
}

// pullReference removes refID from the array field of every user document
// that contains it, bumps those documents' versions, and returns their IDs.
func (s *store[T, PT]) pullReference(ctx context.Context, userID primitive.ObjectID, field string, refID primitive.ObjectID) ([]primitive.ObjectID, error) {
	filter := bson.M{"user_id": userID, field: refID}

	// Find which documents will change, so callers can notify other devices
	cursor, err := s.collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	var rows []struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	ids := make([]primitive.ObjectID, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}

	update := bson.M{
		"$pull": bson.M{field: refID},
		"$inc":  bson.M{"version": 1},
		"$set":  bson.M{"updated_at": time.Now().UTC()},
	}
	if _, err := s.collection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}, "user_id": userID}, update); err != nil {
		return nil, err
	}
	return ids, nil
}

// uniqueIDs returns ids without duplicates, keeping their order.
func uniqueIDs(ids []primitive.ObjectID) []primitive.ObjectID {
	seen := make(map[primitive.ObjectID]struct{}, len(ids))
	unique := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	return unique
}
