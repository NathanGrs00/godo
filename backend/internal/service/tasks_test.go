package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/realtime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskCreateFillsMetaAndPublishes(t *testing.T) {
	f := newFixture()
	in := validTask()
	in.Title = "  Write report  "

	task, err := f.tasks.Create(ctx, f.user, in)
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != "Write report" {
		t.Errorf("title = %q, want trimmed", task.Title)
	}
	if task.Version != 1 || task.UserID != f.user || !task.CreatedAt.Equal(f.now) {
		t.Errorf("meta = %+v", task.Meta)
	}
	if task.TagIDs == nil || task.DeadlineIDs == nil {
		t.Errorf("link slices must be empty, not nil")
	}
	if got := f.events.types(); !equal(got, []string{realtime.TaskCreated}) {
		t.Errorf("events = %v", got)
	}
}

func TestTaskValidation(t *testing.T) {
	f := newFixture()
	cases := []struct {
		name  string
		edit  func(*TaskInput)
		field string
	}{
		{"empty title", func(in *TaskInput) { in.Title = "   " }, "title"},
		{"long title", func(in *TaskInput) { in.Title = strings.Repeat("a", maxTaskTitle+1) }, "title"},
		{"long description", func(in *TaskInput) { in.Description = strings.Repeat("a", maxTaskDescription+1) }, "description"},
		{"missing priority", func(in *TaskInput) { in.Priority = "" }, "priority"},
		{"bad priority", func(in *TaskInput) { in.Priority = "urgent" }, "priority"},
		{"unknown tag", func(in *TaskInput) { in.TagIDs = []primitive.ObjectID{primitive.NewObjectID()} }, "tag_ids"},
		{"unknown deadline", func(in *TaskInput) { in.DeadlineIDs = []primitive.ObjectID{primitive.NewObjectID()} }, "deadline_ids"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := validTask()
			tc.edit(&in)
			_, err := f.tasks.Create(ctx, f.user, in)
			assertValidation(t, err, tc.field)
		})
	}
	if n := len(f.events.types()); n != 0 {
		t.Errorf("invalid input published %d events", n)
	}
}

func TestTaskLinksTagsAndDeadlines(t *testing.T) {
	f := newFixture()
	tag, err := f.tags.Create(ctx, f.user, TagInput{Title: "work", Color: "#8cfbde"})
	if err != nil {
		t.Fatal(err)
	}
	deadline, err := f.deadlines.Create(ctx, f.user, DeadlineInput{Title: "Launch", Date: f.now.AddDate(0, 0, 7)})
	if err != nil {
		t.Fatal(err)
	}

	in := validTask()
	in.TagIDs = []primitive.ObjectID{tag.ID, tag.ID}
	in.DeadlineIDs = []primitive.ObjectID{deadline.ID}
	task, err := f.tasks.Create(ctx, f.user, in)
	if err != nil {
		t.Fatal(err)
	}
	if len(task.TagIDs) != 1 || len(task.DeadlineIDs) != 1 {
		t.Fatalf("links = %v / %v, want one each (duplicates removed)", task.TagIDs, task.DeadlineIDs)
	}

	// Another user cannot link to these
	_, err = f.tasks.Create(ctx, primitive.NewObjectID(), in)
	assertValidation(t, err, "deadline_ids")
}

func TestTaskUpdate(t *testing.T) {
	f := newFixture()
	task, err := f.tasks.Create(ctx, f.user, validTask())
	if err != nil {
		t.Fatal(err)
	}

	in := validTask()
	in.Completed = true
	in.Version = task.Version
	updated, err := f.tasks.Update(ctx, f.user, task.ID, in)
	if err != nil {
		t.Fatal(err)
	}
	if !updated.Completed || updated.Version != 2 || !updated.CreatedAt.Equal(task.CreatedAt) {
		t.Fatalf("updated = %+v", updated)
	}

	t.Run("stale version conflicts", func(t *testing.T) {
		in.Version = 1
		if _, err := f.tasks.Update(ctx, f.user, task.ID, in); !errors.Is(err, ErrConflict) {
			t.Fatalf("err = %v, want ErrConflict", err)
		}
	})
	t.Run("missing version is rejected", func(t *testing.T) {
		in.Version = 0
		_, err := f.tasks.Update(ctx, f.user, task.ID, in)
		assertValidation(t, err, "version")
	})
	t.Run("other user gets not found", func(t *testing.T) {
		in.Version = 2
		if _, err := f.tasks.Update(ctx, primitive.NewObjectID(), task.ID, in); !errors.Is(err, ErrNotFound) {
			t.Fatalf("err = %v, want ErrNotFound", err)
		}
	})

	want := []string{realtime.TaskCreated, realtime.TaskUpdated}
	if got := f.events.types(); !equal(got, want) {
		t.Errorf("events = %v, want %v", got, want)
	}
}

func TestTaskDelete(t *testing.T) {
	f := newFixture()
	task, err := f.tasks.Create(ctx, f.user, validTask())
	if err != nil {
		t.Fatal(err)
	}
	if err := f.tasks.Delete(ctx, primitive.NewObjectID(), task.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete by other user err = %v", err)
	}
	if err := f.tasks.Delete(ctx, f.user, task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := f.tasks.Get(ctx, f.user, task.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get after delete err = %v", err)
	}
	last := f.events.events[len(f.events.events)-1]
	if last.Type != realtime.TaskDeleted || last.ID != task.ID.Hex() || last.Data != nil {
		t.Fatalf("last event = %+v", last)
	}
}

func TestPriorityValid(t *testing.T) {
	for _, p := range []models.Priority{models.PriorityLow, models.PriorityMedium, models.PriorityHigh} {
		if !p.Valid() {
			t.Errorf("%q should be valid", p)
		}
	}
	if models.Priority("High").Valid() {
		t.Error("priorities are case sensitive")
	}
}
