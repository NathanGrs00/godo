package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is an account. All tasks, deadlines and tags belong to a user,
// which is what lets the same data show up on every device they log in on.
type User struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Email string             `bson:"email" json:"email"`
	// PasswordHash is a bcrypt hash and is never sent to clients.
	PasswordHash string    `bson:"password_hash" json:"-"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}
