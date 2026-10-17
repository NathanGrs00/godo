package mongo

import (
	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

// tagRepo is the MongoDB implementation of repository.TagRepository.
type tagRepo struct {
	store[models.Tag, *models.Tag]
}

// NewTagRepository creates a TagRepository backed by the "tags" collection.
func NewTagRepository(db *mongo.Database) repository.TagRepository {
	return &tagRepo{store[models.Tag, *models.Tag]{collection: db.Collection(tagsCollection)}}
}
