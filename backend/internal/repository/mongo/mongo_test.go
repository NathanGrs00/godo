package mongo_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	mongorepo "github.com/NathanGrs00/godo/internal/repository/mongo"
	"github.com/NathanGrs00/godo/internal/repository/repotest"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestContract runs the repository suite against a real MongoDB.
// It is skipped unless MONGO_TEST_URI is set, e.g.
//
//	MONGO_TEST_URI=mongodb://localhost:27017 go test ./...
func TestContract(t *testing.T) {
	uri := os.Getenv("MONGO_TEST_URI")
	if uri == "" {
		t.Skip("MONGO_TEST_URI not set; skipping MongoDB integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { client.Disconnect(context.Background()) })

	repotest.Run(t, func(t *testing.T) repotest.Repositories {
		// Each sub-test gets its own throwaway database
		db := client.Database(fmt.Sprintf("godo_test_%s", primitive.NewObjectID().Hex()))
		t.Cleanup(func() { db.Drop(context.Background()) })
		if err := mongorepo.EnsureIndexes(context.Background(), db); err != nil {
			t.Fatal(err)
		}
		return repotest.Repositories{
			Tasks:     mongorepo.NewTaskRepository(db),
			Deadlines: mongorepo.NewDeadlineRepository(db),
			Tags:      mongorepo.NewTagRepository(db),
			Users:     mongorepo.NewUserRepository(db),
		}
	})
}
