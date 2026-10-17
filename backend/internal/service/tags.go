package service

import (
	"context"
	"regexp"
	"strings"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/realtime"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// maxTagTitle is the longest tag name allowed.
const maxTagTitle = 50

// hexColor matches colours like "#8CFBDE".
var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// TagInput is what a client sends to create or update a tag.
type TagInput struct {
	Title string `json:"title"`
	Color string `json:"color"`
	// Version is the version the client last saw. Required for updates.
	Version int64 `json:"version"`
}

// TagService manages a user's tags.
type TagService struct {
	tags      repository.TagRepository
	tasks     *TaskService
	deadlines *DeadlineService
	publisher Publisher
	now       clock
}

// NewTagService creates a TagService. tasks and deadlines are used to unlink
// and re-announce items when a tag is deleted.
func NewTagService(tags repository.TagRepository, tasks *TaskService, deadlines *DeadlineService, publisher Publisher) *TagService {
	return &TagService{tags: tags, tasks: tasks, deadlines: deadlines, publisher: publisher, now: systemClock}
}

// List returns all of the user's tags.
func (s *TagService) List(ctx context.Context, userID primitive.ObjectID) ([]*models.Tag, error) {
	return s.tags.List(ctx, userID)
}

// Get returns one of the user's tags.
func (s *TagService) Get(ctx context.Context, userID, id primitive.ObjectID) (*models.Tag, error) {
	return s.tags.Get(ctx, userID, id)
}

// Create validates and stores a new tag.
func (s *TagService) Create(ctx context.Context, userID primitive.ObjectID, in TagInput) (*models.Tag, error) {
	tag := &models.Tag{Meta: newMeta(userID, s.now())}
	if err := applyTag(tag, in); err != nil {
		return nil, err
	}
	if err := s.tags.Create(ctx, tag); err != nil {
		return nil, err
	}
	s.publish(userID, realtime.TagCreated, tag)
	return tag, nil
}

// Update validates and replaces a tag. It returns ErrConflict when the tag
// changed since in.Version.
func (s *TagService) Update(ctx context.Context, userID, id primitive.ObjectID, in TagInput) (*models.Tag, error) {
	if err := checkVersion(in.Version); err != nil {
		return nil, err
	}
	tag, err := s.tags.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if err := applyTag(tag, in); err != nil {
		return nil, err
	}
	tag.Version = in.Version
	tag.UpdatedAt = s.now()
	if err := s.tags.Update(ctx, tag); err != nil {
		return nil, err
	}
	s.publish(userID, realtime.TagUpdated, tag)
	return tag, nil
}

// Delete removes a tag and unlinks it from every task and deadline.
func (s *TagService) Delete(ctx context.Context, userID, id primitive.ObjectID) error {
	if err := s.tags.Delete(ctx, userID, id); err != nil {
		return err
	}
	s.publisher.Publish(userID, realtime.Event{Type: realtime.TagDeleted, ID: id.Hex()})

	changedTasks, err := s.tasks.tasks.RemoveTag(ctx, userID, id)
	if err != nil {
		return err
	}
	s.tasks.publishChanged(ctx, userID, changedTasks)

	changedDeadlines, err := s.deadlines.deadlines.RemoveTag(ctx, userID, id)
	if err != nil {
		return err
	}
	s.deadlines.publishChanged(ctx, userID, changedDeadlines)
	return nil
}

// applyTag validates in and copies it onto tag.
func applyTag(tag *models.Tag, in TagInput) error {
	title, err := cleanText("title", in.Title, true, maxTagTitle)
	if err != nil {
		return err
	}
	if !hexColor.MatchString(in.Color) {
		return invalid("color", "must be a hex colour like #8CFBDE")
	}
	tag.Title = title
	tag.Color = strings.ToUpper(in.Color)
	return nil
}

// publish sends a created/updated event for tag.
func (s *TagService) publish(userID primitive.ObjectID, eventType string, tag *models.Tag) {
	s.publisher.Publish(userID, realtime.Event{Type: eventType, ID: tag.ID.Hex(), Data: tag})
}
