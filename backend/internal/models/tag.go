package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	ID     primitive.ObjectID `json:"id"`
	Name   string             `json:"name"`
	Color  string             `json:"color"`
	UserId primitive.ObjectID `json:"user_id"`
}
