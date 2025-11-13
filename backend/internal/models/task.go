package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Priority    string               `json:"priority"`
	Completed   bool                 `json:"completed"`
	UserId      primitive.ObjectID   `json:"user_id"`
	DeadlineIDs []primitive.ObjectID `json:"deadline_ids"`
	TagIDs      []primitive.ObjectID `json:"tag_ids"`
}
