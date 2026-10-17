package service

import (
	"testing"

	"github.com/NathanGrs00/godo/internal/realtime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTagValidation(t *testing.T) {
	f := newFixture()
	for _, color := range []string{"", "red", "#FFF", "#GGGGGG", "8CFBDE"} {
		_, err := f.tags.Create(ctx, f.user, TagInput{Title: "x", Color: color})
		assertValidation(t, err, "color")
	}
	_, err := f.tags.Create(ctx, f.user, TagInput{Color: "#000000"})
	assertValidation(t, err, "title")

	tag, err := f.tags.Create(ctx, f.user, TagInput{Title: "ok", Color: "#8cfbde"})
	if err != nil {
		t.Fatal(err)
	}
	if tag.Color != "#8CFBDE" {
		t.Errorf("color = %q, want upper case", tag.Color)
	}
}

func TestTagUpdate(t *testing.T) {
	f := newFixture()
	tag, err := f.tags.Create(ctx, f.user, TagInput{Title: "old", Color: "#000000"})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := f.tags.Update(ctx, f.user, tag.ID, TagInput{Title: "new", Color: "#FFFFFF", Version: 1})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "new" || updated.Version != 2 {
		t.Fatalf("updated = %+v", updated)
	}
}

func TestTagDeleteUnlinksEverywhere(t *testing.T) {
	f := newFixture()
	tag, err := f.tags.Create(ctx, f.user, TagInput{Title: "work", Color: "#000000"})
	if err != nil {
		t.Fatal(err)
	}
	ids := []primitive.ObjectID{tag.ID}

	in := validTask()
	in.TagIDs = ids
	task, err := f.tasks.Create(ctx, f.user, in)
	if err != nil {
		t.Fatal(err)
	}
	deadline, err := f.deadlines.Create(ctx, f.user, DeadlineInput{Title: "d", Date: f.now, TagIDs: ids})
	if err != nil {
		t.Fatal(err)
	}

	if err := f.tags.Delete(ctx, f.user, tag.ID); err != nil {
		t.Fatal(err)
	}

	gotTask, _ := f.tasks.Get(ctx, f.user, task.ID)
	gotDeadline, _ := f.deadlines.Get(ctx, f.user, deadline.ID)
	if len(gotTask.TagIDs) != 0 || len(gotDeadline.TagIDs) != 0 {
		t.Fatalf("tag still linked: task %v deadline %v", gotTask.TagIDs, gotDeadline.TagIDs)
	}

	want := []string{realtime.TagCreated, realtime.TaskCreated, realtime.DeadlineCreated, realtime.TagDeleted, realtime.TaskUpdated, realtime.DeadlineUpdated}
	if got := f.events.types(); !equal(got, want) {
		t.Fatalf("events = %v, want %v", got, want)
	}
}
