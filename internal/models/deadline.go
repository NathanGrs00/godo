package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Deadline struct {
	ID          primitive.ObjectID
	Title       string
	Description string
	Date        time.Time
	Passed      bool
	UserId      primitive.ObjectID
	TaskIDs     []primitive.ObjectID
	TagIDs      []primitive.ObjectID
}
