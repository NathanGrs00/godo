package mongo

import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// deadlineRepo is the MongoDB implementation of repository.DeadlineRepository.
type deadlineRepo struct {
	store[models.Deadline, *models.Deadline]
}

// NewDeadlineRepository creates a DeadlineRepository backed by the "deadlines" collection.
func NewDeadlineRepository(db *mongo.Database) repository.DeadlineRepository {
	return &deadlineRepo{store[models.Deadline, *models.Deadline]{collection: db.Collection(deadlinesCollection)}}
}

// RemoveTag unlinks a tag from every deadline of the user.
func (r *deadlineRepo) RemoveTag(ctx context.Context, userID, tagID primitive.ObjectID) ([]primitive.ObjectID, error) {
	return r.pullReference(ctx, userID, "tag_ids", tagID)
}
