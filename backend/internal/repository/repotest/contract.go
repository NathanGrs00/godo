// Package repotest holds a test suite every repository implementation must
// pass. Running the same suite against the in-memory and MongoDB versions
// guarantees they behave the same, so tests using the in-memory version
// are trustworthy.
package repotest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repositories is one full set of repository implementations.
type Repositories struct {
	Tasks     repository.TaskRepository
	Deadlines repository.DeadlineRepository
	Tags      repository.TagRepository
	Users     repository.UserRepository
}

// Run runs the contract suite. newRepos must return fresh, empty repositories.
func Run(t *testing.T, newRepos func(t *testing.T) Repositories) {
	ctx := context.Background()

	t.Run("create get and list are scoped to the owner", func(t *testing.T) {
		repos := newRepos(t)
		alice, bob := primitive.NewObjectID(), primitive.NewObjectID()
		first := newTask(alice, "first", time.Unix(1, 0))
		second := newTask(alice, "second", time.Unix(2, 0))
		other := newTask(bob, "bob's", time.Unix(3, 0))
		for _, task := range []*models.Task{first, second, other} {
			must(t, repos.Tasks.Create(ctx, task))
		}

		got, err := repos.Tasks.Get(ctx, alice, first.ID)
		must(t, err)
		if got.Title != "first" {
			t.Fatalf("Get title = %q, want first", got.Title)
		}

		if _, err := repos.Tasks.Get(ctx, bob, first.ID); !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("Get by other user err = %v, want ErrNotFound", err)
		}

		list, err := repos.Tasks.List(ctx, alice)
		must(t, err)
		if len(list) != 2 || list[0].Title != "first" || list[1].Title != "second" {
			t.Fatalf("List = %v, want [first second]", titles(list))
		}

		empty, err := repos.Tasks.List(ctx, primitive.NewObjectID())
		must(t, err)
		if empty == nil || len(empty) != 0 {
			t.Fatalf("List for unknown user = %#v, want empty non-nil slice", empty)
		}
	})

	t.Run("update checks and bumps the version", func(t *testing.T) {
		repos := newRepos(t)
		user := primitive.NewObjectID()
		task := newTask(user, "v1", time.Now())
		must(t, repos.Tasks.Create(ctx, task))

		task.Title = "v2"
		must(t, repos.Tasks.Update(ctx, task))
		if task.Version != 2 {
			t.Fatalf("Version after update = %d, want 2", task.Version)
		}

		stale := *task
		stale.Version = 1
		stale.Title = "stale"
		if err := repos.Tasks.Update(ctx, &stale); !errors.Is(err, repository.ErrConflict) {
			t.Fatalf("stale Update err = %v, want ErrConflict", err)
		}
		if stale.Version != 1 {
			t.Fatalf("failed Update changed Version to %d", stale.Version)
		}

		got, err := repos.Tasks.Get(ctx, user, task.ID)
		must(t, err)
		if got.Title != "v2" || got.Version != 2 {
			t.Fatalf("stored = %q v%d, want v2 v2", got.Title, got.Version)
		}

		missing := newTask(user, "missing", time.Now())
		if err := repos.Tasks.Update(ctx, missing); !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("Update missing err = %v, want ErrNotFound", err)
		}

		foreign := *got
		foreign.UserID = primitive.NewObjectID()
		if err := repos.Tasks.Update(ctx, &foreign); !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("Update by other user err = %v, want ErrNotFound", err)
		}
	})

	t.Run("delete is scoped to the owner", func(t *testing.T) {
		repos := newRepos(t)
		user := primitive.NewObjectID()
		tag := &models.Tag{Meta: newMeta(user, time.Now()), Title: "x", Color: "#000000"}
		must(t, repos.Tags.Create(ctx, tag))

		if err := repos.Tags.Delete(ctx, primitive.NewObjectID(), tag.ID); !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("Delete by other user err = %v, want ErrNotFound", err)
		}
		must(t, repos.Tags.Delete(ctx, user, tag.ID))
		if err := repos.Tags.Delete(ctx, user, tag.ID); !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("second Delete err = %v, want ErrNotFound", err)
		}
	})

	t.Run("all exist", func(t *testing.T) {
		repos := newRepos(t)
		user := primitive.NewObjectID()
		tag := &models.Tag{Meta: newMeta(user, time.Now()), Title: "x", Color: "#000000"}
		must(t, repos.Tags.Create(ctx, tag))

		cases := []struct {
			name string
			user primitive.ObjectID
			ids  []primitive.ObjectID
			want bool
		}{
			{"none", user, nil, true},
			{"own", user, []primitive.ObjectID{tag.ID}, true},
			{"duplicates", user, []primitive.ObjectID{tag.ID, tag.ID}, true},
			{"unknown", user, []primitive.ObjectID{tag.ID, primitive.NewObjectID()}, false},
			{"other user", primitive.NewObjectID(), []primitive.ObjectID{tag.ID}, false},
		}
		for _, tc := range cases {
			got, err := repos.Tags.AllExist(ctx, tc.user, tc.ids)
			must(t, err)
			if got != tc.want {
				t.Errorf("%s: AllExist = %v, want %v", tc.name, got, tc.want)
			}
		}
	})

	t.Run("remove references", func(t *testing.T) {
		repos := newRepos(t)
		user := primitive.NewObjectID()
		tagID, keepID, deadlineID := primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID()

		linked := newTask(user, "linked", time.Now())
		linked.TagIDs = []primitive.ObjectID{tagID, keepID}
		linked.DeadlineIDs = []primitive.ObjectID{deadlineID}
		unlinked := newTask(user, "unlinked", time.Now())
		otherUser := newTask(primitive.NewObjectID(), "other", time.Now())
		otherUser.TagIDs = []primitive.ObjectID{tagID}
		for _, task := range []*models.Task{linked, unlinked, otherUser} {
			must(t, repos.Tasks.Create(ctx, task))
		}

		changed, err := repos.Tasks.RemoveTag(ctx, user, tagID)
		must(t, err)
		if len(changed) != 1 || changed[0] != linked.ID {
			t.Fatalf("RemoveTag changed = %v, want [%v]", changed, linked.ID)
		}
		got, err := repos.Tasks.Get(ctx, user, linked.ID)
		must(t, err)
		if len(got.TagIDs) != 1 || got.TagIDs[0] != keepID || got.Version != 2 {
			t.Fatalf("after RemoveTag: tags %v version %d", got.TagIDs, got.Version)
		}
		other, err := repos.Tasks.Get(ctx, otherUser.UserID, otherUser.ID)
		must(t, err)
		if len(other.TagIDs) != 1 {
			t.Fatalf("RemoveTag touched another user's task")
		}

		changed, err = repos.Tasks.RemoveDeadline(ctx, user, deadlineID)
		must(t, err)
		if len(changed) != 1 {
			t.Fatalf("RemoveDeadline changed %d tasks, want 1", len(changed))
		}

		deadline := &models.Deadline{Meta: newMeta(user, time.Now()), Title: "d", Date: time.Now().UTC().Truncate(time.Millisecond), TagIDs: []primitive.ObjectID{tagID}}
		must(t, repos.Deadlines.Create(ctx, deadline))
		changed, err = repos.Deadlines.RemoveTag(ctx, user, tagID)
		must(t, err)
		if len(changed) != 1 {
			t.Fatalf("Deadlines.RemoveTag changed %d, want 1", len(changed))
		}
	})

	t.Run("users", func(t *testing.T) {
		repos := newRepos(t)
		user := &models.User{ID: primitive.NewObjectID(), Email: "a@example.com", PasswordHash: "hash", CreatedAt: time.Now().UTC().Truncate(time.Millisecond)}
		must(t, repos.Users.Create(ctx, user))

		dup := &models.User{ID: primitive.NewObjectID(), Email: "a@example.com"}
		if err := repos.Users.Create(ctx, dup); !errors.Is(err, repository.ErrDuplicate) {
			t.Fatalf("duplicate email err = %v, want ErrDuplicate", err)
		}

		byEmail, err := repos.Users.GetByEmail(ctx, "a@example.com")
		must(t, err)
		byID, err := repos.Users.GetByID(ctx, user.ID)
		must(t, err)
		if byEmail.ID != user.ID || byID.Email != user.Email || byID.PasswordHash != "hash" {
			t.Fatalf("lookups returned %+v and %+v", byEmail, byID)
		}
		if _, err := repos.Users.GetByEmail(ctx, "nobody@example.com"); !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("unknown email err = %v, want ErrNotFound", err)
		}
	})
}

func newMeta(user primitive.ObjectID, created time.Time) models.Meta {
	created = created.UTC().Truncate(time.Millisecond)
	return models.Meta{ID: primitive.NewObjectID(), UserID: user, Version: 1, CreatedAt: created, UpdatedAt: created}
}

func newTask(user primitive.ObjectID, title string, created time.Time) *models.Task {
	return &models.Task{
		Meta:        newMeta(user, created),
		Title:       title,
		Priority:    models.PriorityLow,
		DeadlineIDs: []primitive.ObjectID{},
		TagIDs:      []primitive.ObjectID{},
	}
}

func titles(tasks []*models.Task) []string {
	out := make([]string, len(tasks))
	for i, task := range tasks {
		out[i] = task.Title
	}
	return out
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
