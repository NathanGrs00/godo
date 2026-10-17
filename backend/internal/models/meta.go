package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Meta holds the fields every user-owned document shares.
// It is embedded (inlined) in Task, Deadline and Tag.
type Meta struct {
	// ID is the document's unique identifier.
	ID primitive.ObjectID `bson:"_id" json:"id"`
	// UserID is the owner of the document. It is never sent to clients.
	UserID primitive.ObjectID `bson:"user_id" json:"-"`
	// Version increases by one on every write. Clients send back the version
	// they last saw, so an edit made on a stale copy is rejected instead of
	// silently overwriting a newer edit from another device.
	Version int64 `bson:"version" json:"version"`
	// CreatedAt is when the document was first stored.
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	// UpdatedAt is when the document was last changed.
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// GetMeta returns a pointer to the shared fields so generic code
// (repositories, services) can read and update them.
func (m *Meta) GetMeta() *Meta { return m }

// Entity is implemented by every user-owned model through its embedded Meta.
type Entity interface {
	GetMeta() *Meta
}
