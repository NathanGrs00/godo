// Package memory implements the repository interfaces in memory.
//
// It behaves like the MongoDB implementation (user scoping, versions,
// not-found and conflict errors) and is used by tests so they run without
// a database.
package memory

import (
	"context"
	"slices"
	"sync"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// store is a generic, user-scoped, thread-safe collection.
type store[T any, PT interface {
	*T
	models.Entity
}] struct {
	mu    sync.RWMutex
	items map[primitive.ObjectID]PT
	// order keeps insertion order so List is deterministic (oldest first)
	order []primitive.ObjectID
}

func newStore[T any, PT interface {
	*T
	models.Entity
}]() store[T, PT] {
	return store[T, PT]{items: map[primitive.ObjectID]PT{}}
}

// clone deep-copies an item through BSON, exactly like a database round trip.
func clone[T any, PT interface {
	*T
	models.Entity
}](item PT) PT {
	raw, err := bson.Marshal(item)
	if err != nil {
		panic(err)
	}
	var out T
	if err := bson.Unmarshal(raw, &out); err != nil {
		panic(err)
	}
	return &out
}

func (s *store[T, PT]) Create(_ context.Context, item PT) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := item.GetMeta().ID
	if _, ok := s.items[id]; ok {
		return repository.ErrDuplicate
	}
	s.items[id] = clone[T, PT](item)
	s.order = append(s.order, id)
	return nil
}

func (s *store[T, PT]) Get(_ context.Context, userID, id primitive.ObjectID) (PT, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok || item.GetMeta().UserID != userID {
		return nil, repository.ErrNotFound
	}
	return clone[T, PT](item), nil
}

func (s *store[T, PT]) List(_ context.Context, userID primitive.ObjectID) ([]PT, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := []PT{}
	for _, id := range s.order {
		if item := s.items[id]; item.GetMeta().UserID == userID {
			items = append(items, clone[T, PT](item))
		}
	}
	return items, nil
}

func (s *store[T, PT]) Update(_ context.Context, item PT) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	meta := item.GetMeta()
	stored, ok := s.items[meta.ID]
	if !ok || stored.GetMeta().UserID != meta.UserID {
		return repository.ErrNotFound
	}
	if stored.GetMeta().Version != meta.Version {
		return repository.ErrConflict
	}
	meta.Version++
	s.items[meta.ID] = clone[T, PT](item)
	return nil
}

func (s *store[T, PT]) Delete(_ context.Context, userID, id primitive.ObjectID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok || item.GetMeta().UserID != userID {
		return repository.ErrNotFound
	}
	delete(s.items, id)
	s.order = slices.DeleteFunc(s.order, func(o primitive.ObjectID) bool { return o == id })
	return nil
}

func (s *store[T, PT]) AllExist(_ context.Context, userID primitive.ObjectID, ids []primitive.ObjectID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, id := range ids {
		item, ok := s.items[id]
		if !ok || item.GetMeta().UserID != userID {
			return false, nil
		}
	}
	return true, nil
}

// pullReference removes refID from the slice that field returns, for every
// item of the user, bumping versions of the items that changed.
func (s *store[T, PT]) pullReference(userID, refID primitive.ObjectID, field func(PT) *[]primitive.ObjectID) []primitive.ObjectID {
	s.mu.Lock()
	defer s.mu.Unlock()
	var changed []primitive.ObjectID
	for _, id := range s.order {
		item := s.items[id]
		if item.GetMeta().UserID != userID {
			continue
		}
		refs := field(item)
		if !slices.Contains(*refs, refID) {
			continue
		}
		*refs = slices.DeleteFunc(*refs, func(r primitive.ObjectID) bool { return r == refID })
		item.GetMeta().Version++
		changed = append(changed, id)
	}
	return changed
}
