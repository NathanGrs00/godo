package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Deadline struct {
	ID          primitive.ObjectID   `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Date        time.Time            `json:"date"`
	Passed      bool                 `json:"passed"`
	UserId      primitive.ObjectID   `json:"user_id"`
	TaskIDs     []primitive.ObjectID `json:"task_ids"`
	TagIDs      []primitive.ObjectID `json:"tag_ids"`
}
