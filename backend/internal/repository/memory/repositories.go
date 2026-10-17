package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskRepository is an in-memory repository.TaskRepository.
type TaskRepository struct {
	store[models.Task, *models.Task]
}

// NewTaskRepository creates an empty TaskRepository.
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{newStore[models.Task, *models.Task]()}
}

// RemoveTag unlinks a tag from every task of the user.
func (r *TaskRepository) RemoveTag(_ context.Context, userID, tagID primitive.ObjectID) ([]primitive.ObjectID, error) {
	return r.pullReference(userID, tagID, func(t *models.Task) *[]primitive.ObjectID { return &t.TagIDs }), nil
}

// RemoveDeadline unlinks a deadline from every task of the user.
func (r *TaskRepository) RemoveDeadline(_ context.Context, userID, deadlineID primitive.ObjectID) ([]primitive.ObjectID, error) {
	return r.pullReference(userID, deadlineID, func(t *models.Task) *[]primitive.ObjectID { return &t.DeadlineIDs }), nil
}

// DeadlineRepository is an in-memory repository.DeadlineRepository.
type DeadlineRepository struct {
	store[models.Deadline, *models.Deadline]
}

// NewDeadlineRepository creates an empty DeadlineRepository.
func NewDeadlineRepository() *DeadlineRepository {
	return &DeadlineRepository{newStore[models.Deadline, *models.Deadline]()}
}

// RemoveTag unlinks a tag from every deadline of the user.
func (r *DeadlineRepository) RemoveTag(_ context.Context, userID, tagID primitive.ObjectID) ([]primitive.ObjectID, error) {
	return r.pullReference(userID, tagID, func(d *models.Deadline) *[]primitive.ObjectID { return &d.TagIDs }), nil
}

// TagRepository is an in-memory repository.TagRepository.
type TagRepository struct {
	store[models.Tag, *models.Tag]
}

// NewTagRepository creates an empty TagRepository.
func NewTagRepository() *TagRepository {
	return &TagRepository{newStore[models.Tag, *models.Tag]()}
}

// UserRepository is an in-memory repository.UserRepository.
type UserRepository struct {
	mu    sync.RWMutex
	users map[primitive.ObjectID]models.User
}

// NewUserRepository creates an empty UserRepository.
func NewUserRepository() *UserRepository {
	return &UserRepository{users: map[primitive.ObjectID]models.User{}}
}

// Create stores a user, rejecting duplicate emails.
func (r *UserRepository) Create(_ context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if strings.EqualFold(u.Email, user.Email) {
			return repository.ErrDuplicate
		}
	}
	r.users[user.ID] = *user
	return nil
}

// GetByID finds a user by ID.
func (r *UserRepository) GetByID(_ context.Context, id primitive.ObjectID) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return &u, nil
}

// GetByEmail finds a user by email.
func (r *UserRepository) GetByEmail(_ context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, repository.ErrNotFound
}

// Compile-time checks that the in-memory types satisfy the interfaces.
var (
	_ repository.TaskRepository     = (*TaskRepository)(nil)
	_ repository.DeadlineRepository = (*DeadlineRepository)(nil)
	_ repository.TagRepository      = (*TagRepository)(nil)
	_ repository.UserRepository     = (*UserRepository)(nil)
)
