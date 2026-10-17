package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection names.
const (
	usersCollection     = "users"
	tasksCollection     = "tasks"
	deadlinesCollection = "deadlines"
	tagsCollection      = "tags"
)

// EnsureIndexes creates the indexes the repositories rely on. It is safe to
// run on every start: existing indexes are left alone.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	// Emails must be unique so two accounts cannot share one
	_, err := db.Collection(usersCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Every list query filters by owner and sorts by creation time
	for _, name := range []string{tasksCollection, deadlinesCollection, tagsCollection} {
		_, err := db.Collection(name).Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: 1}},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
