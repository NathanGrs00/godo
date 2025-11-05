package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID
	Title       string
	Description string
	Priority    string
	Completed   bool
	UserId      primitive.ObjectID
	DeadlineIDs []primitive.ObjectID
	TagIDs      []primitive.ObjectID
}
