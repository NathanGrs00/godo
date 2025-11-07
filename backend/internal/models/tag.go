package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	ID     primitive.ObjectID `json:"id"`
	Title  string             `json:"title"`
	Color  string             `json:"color"`
	UserId primitive.ObjectID `json:"user_id"`
}
