package mongo

import (
	"context"

	"github.com/NathanGrs00/godo/internal/models"
	_interface "github.com/NathanGrs00/godo/internal/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

// taskRepo represents the MongoDB implementation of the Task repository.
type taskRepo struct {
	collection *mongo.Collection
}

// NewTaskRepository creates a new instance of TaskRepository using the provided MongoDB database.
func NewTaskRepository(db *mongo.Database) _interface.TaskRepository {
	return &taskRepo{collection: db.Collection("tasks")}
}

// Create adds a new task to the MongoDB collection.
// r is like self in Python, it refers to the instance of taskRepo.
func (r *taskRepo) Create(ctx context.Context, task *models.Task) error {
	_, err := r.collection.InsertOne(ctx, task)
	return err
}

// GetById retrieves a task by its ID from the MongoDB collection.
func (r *taskRepo) GetById(ctx context.Context, id primitive.ObjectID) (*models.Task, error) {
	// task will hold the result of the query
	var task models.Task
	// Filter to find the task by its ID
	filter := bson.M{"_id": id}
	// Perform the query and decode the result into the task variable
	err := r.collection.FindOne(ctx, filter).Decode(&task)
	// Return nil for the task and the error if there was an issue
	if err != nil {
		return nil, err
	}
	// Return the found task
	return &task, nil
}

// GetByUserId retrieves all tasks for a specific user from the MongoDB collection.
func (r *taskRepo) GetByUserId(ctx context.Context, userID primitive.ObjectID) ([]*models.Task, error) {
	// Filter to find tasks by user ID
	filter := bson.M{"user_id": userID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed after processing
	defer cursor.Close(ctx)

	// tasks will hold the list of tasks for the user
	var tasks []*models.Task
	// Iterate through the cursor
	for cursor.Next(ctx) {
		// Decode each document into a Task struct
		var task models.Task
		// If error occurs during decoding, return it
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		// Append the task to the tasks slice
		tasks = append(tasks, &task)
	}

	// Return the list of tasks
	return tasks, nil
}

// Update modifies an existing task in the MongoDB collection.
func (r *taskRepo) Update(ctx context.Context, task *models.Task) error {
	// Filter to find the task by its ID
	filter := bson.M{"_id": task.ID}
	// Update the task document with the new values
	update := bson.D{{"$set", task}}
	// Perform the update operation
	_, err := r.collection.UpdateOne(ctx, filter, update)
	// Return the error if any
	return err
}

// Delete removes an existing task in the MongoDB collection.
func (r *taskRepo) Delete(ctx context.Context, taskID primitive.ObjectID) error {
	// Filter to find the task by ID
	filter := bson.M{"_id": taskID}
	// Perform the update operation
	_, err := r.collection.DeleteOne(ctx, filter)
	// Return the error if any
	return err
}

// GetAll to get all tasks
func (r *taskRepo) GetAll(ctx context.Context) ([]*models.Task, error) {
	// Make a slice to hold tasks
	var tasks []*models.Task
	// Find all tasks in the collection
	cursor, err := r.collection.Find(ctx, bson.D{})

	// Check for errors during the find operation
	if err != nil {
		// Return nil and the error if any
		return nil, err
	}
	// Ensure the cursor is closed after processing
	defer cursor.Close(ctx)
	// Iterate through the cursor to decode each task
	for cursor.Next(ctx) {
		// Decode the current document into a Task struct
		var task models.Task
		// Check for errors during decoding
		if err := cursor.Decode(&task); err != nil {
			// Return nil and the error if any
			return nil, err
		}
		// Append the decoded task to the tasks slice
		tasks = append(tasks, &task)
	}
	// Return the slice of tasks and nil for error
	return tasks, nil
}
