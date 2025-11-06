package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID   `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Priority    string               `json:"priority"`
	Completed   bool                 `json:"completed"`
	UserId      primitive.ObjectID   `json:"userId"`
	DeadlineIDs []primitive.ObjectID `json:"deadlineIds"`
	TagIDs      []primitive.ObjectID `json:"tagIds"`
}
