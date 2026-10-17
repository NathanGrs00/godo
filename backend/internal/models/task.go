package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Priority is how urgent a task is.
type Priority string

// The allowed task priorities.
const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Valid reports whether p is one of the allowed priorities.
func (p Priority) Valid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}

// Task is a single to-do item.
type Task struct {
	Meta        `bson:",inline"`
	Title       string   `bson:"title" json:"title"`
	Description string   `bson:"description" json:"description"`
	Priority    Priority `bson:"priority" json:"priority"`
	Completed   bool     `bson:"completed" json:"completed"`
	// DeadlineIDs links the task to deadlines. The task owns this link:
	// a deadline's tasks are all tasks whose DeadlineIDs contain it.
	DeadlineIDs []primitive.ObjectID `bson:"deadline_ids" json:"deadline_ids"`
	TagIDs      []primitive.ObjectID `bson:"tag_ids" json:"tag_ids"`
}
