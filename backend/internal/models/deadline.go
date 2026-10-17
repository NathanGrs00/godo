package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Deadline is a dated milestone that tasks can be linked to.
type Deadline struct {
	Meta        `bson:",inline"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	Date        time.Time            `bson:"date" json:"date"`
	TagIDs      []primitive.ObjectID `bson:"tag_ids" json:"tag_ids"`
	// Passed is computed when the deadline is read and is not stored.
	Passed bool `bson:"-" json:"passed"`
}
