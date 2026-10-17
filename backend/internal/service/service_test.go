package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/NathanGrs00/godo/internal/auth"
	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/realtime"
	"github.com/NathanGrs00/godo/internal/repository/memory"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// recorder is a Publisher that remembers every event.
type recorder struct {
	mu     sync.Mutex
	events []realtime.Event
}

func (r *recorder) Publish(_ primitive.ObjectID, e realtime.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, e)
}

// types returns the event types in order.
func (r *recorder) types() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.events))
	for i, e := range r.events {
		out[i] = e.Type
	}
	return out
}

// fixture is a full set of services on in-memory repositories.
type fixture struct {
	tasks     *TaskService
	deadlines *DeadlineService
	tags      *TagService
	auth      *AuthService
	events    *recorder
	user      primitive.ObjectID
	now       time.Time
}

func newFixture() *fixture {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	fixed := func() time.Time { return now }

	events := &recorder{}
	taskRepo, deadlineRepo, tagRepo := memory.NewTaskRepository(), memory.NewDeadlineRepository(), memory.NewTagRepository()

	tasks := NewTaskService(taskRepo, deadlineRepo, tagRepo, events)
	deadlines := NewDeadlineService(deadlineRepo, tagRepo, tasks, events)
	tags := NewTagService(tagRepo, tasks, deadlines, events)
	authSvc := NewAuthService(memory.NewUserRepository(), auth.NewTokenManager("test-secret-test-secret-test-secret", time.Hour))
	tasks.now, deadlines.now, tags.now, authSvc.now = fixed, fixed, fixed, fixed

	return &fixture{tasks: tasks, deadlines: deadlines, tags: tags, auth: authSvc, events: events, user: primitive.NewObjectID(), now: now}
}

var ctx = context.Background()

func validTask() TaskInput {
	return TaskInput{Title: "Write report", Priority: models.PriorityHigh}
}

func assertValidation(t *testing.T, err error, field string) {
	t.Helper()
	var v *ValidationError
	if !errors.As(err, &v) {
		t.Fatalf("err = %v, want ValidationError on %q", err, field)
	}
	if v.Field != field {
		t.Fatalf("validation field = %q, want %q", v.Field, field)
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
