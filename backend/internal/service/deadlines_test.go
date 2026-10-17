package service

import (
	"testing"
	"time"

	"github.com/NathanGrs00/godo/internal/realtime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeadlinePassedIsComputed(t *testing.T) {
	f := newFixture()
	past, err := f.deadlines.Create(ctx, f.user, DeadlineInput{Title: "Past", Date: f.now.Add(-time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	future, err := f.deadlines.Create(ctx, f.user, DeadlineInput{Title: "Future", Date: f.now.Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if !past.Passed || future.Passed {
		t.Fatalf("passed = %v/%v, want true/false", past.Passed, future.Passed)
	}

	// Time moves on: the future deadline is now passed too
	f.deadlines.now = func() time.Time { return f.now.Add(2 * time.Hour) }
	list, err := f.deadlines.List(ctx, f.user)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range list {
		if !d.Passed {
			t.Errorf("%s should be passed", d.Title)
		}
	}
}

func TestDeadlineValidation(t *testing.T) {
	f := newFixture()
	_, err := f.deadlines.Create(ctx, f.user, DeadlineInput{Title: "No date"})
	assertValidation(t, err, "date")
	_, err = f.deadlines.Create(ctx, f.user, DeadlineInput{Date: f.now})
	assertValidation(t, err, "title")
	_, err = f.deadlines.Create(ctx, f.user, DeadlineInput{Title: "x", Date: f.now, TagIDs: []primitive.ObjectID{primitive.NewObjectID()}})
	assertValidation(t, err, "tag_ids")
}

func TestDeadlineDateIsStoredInUTC(t *testing.T) {
	f := newFixture()
	amsterdam := time.FixedZone("CET", 3600)
	d, err := f.deadlines.Create(ctx, f.user, DeadlineInput{Title: "x", Date: time.Date(2026, 3, 1, 9, 0, 0, 0, amsterdam)})
	if err != nil {
		t.Fatal(err)
	}
	if d.Date.Location() != time.UTC || d.Date.Hour() != 8 {
		t.Fatalf("date = %v, want 08:00 UTC", d.Date)
	}
}

func TestDeadlineDeleteUnlinksTasks(t *testing.T) {
	f := newFixture()
	deadline, err := f.deadlines.Create(ctx, f.user, DeadlineInput{Title: "Launch", Date: f.now})
	if err != nil {
		t.Fatal(err)
	}
	in := validTask()
	in.DeadlineIDs = []primitive.ObjectID{deadline.ID}
	task, err := f.tasks.Create(ctx, f.user, in)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.deadlines.Delete(ctx, f.user, deadline.ID); err != nil {
		t.Fatal(err)
	}

	got, err := f.tasks.Get(ctx, f.user, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.DeadlineIDs) != 0 || got.Version != 2 {
		t.Fatalf("task after deadline delete: links %v version %d", got.DeadlineIDs, got.Version)
	}

	// Other devices learn about both the deletion and the changed task
	want := []string{realtime.DeadlineCreated, realtime.TaskCreated, realtime.DeadlineDeleted, realtime.TaskUpdated}
	if got := f.events.types(); !equal(got, want) {
		t.Fatalf("events = %v, want %v", got, want)
	}
}
