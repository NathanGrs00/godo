// Package realtime fans out change events to every connected device of a user.
//
// When a device changes a task, deadline or tag, the service publishes an
// event here. Every other device of the same user that has the app open is
// subscribed (through the /events stream) and applies the change at once.
package realtime

import (
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event types sent to clients. The entity is the part before the dot.
const (
	TaskCreated     = "task.created"
	TaskUpdated     = "task.updated"
	TaskDeleted     = "task.deleted"
	DeadlineCreated = "deadline.created"
	DeadlineUpdated = "deadline.updated"
	DeadlineDeleted = "deadline.deleted"
	TagCreated      = "tag.created"
	TagUpdated      = "tag.updated"
	TagDeleted      = "tag.deleted"
)

// Event describes one change.
type Event struct {
	// Type is one of the constants above.
	Type string `json:"type"`
	// ID is the ID of the changed document.
	ID string `json:"id"`
	// Data is the document after the change. It is empty for deletions.
	Data any `json:"data,omitempty"`
}

// bufferSize is how many events a subscriber may fall behind before it is
// disconnected. A disconnected client reconnects and reloads everything, so
// no change is ever silently lost.
const bufferSize = 64

// subscriber is one open event stream.
type subscriber struct {
	events chan Event
}

// Hub keeps track of open streams per user. It is safe for concurrent use.
//
// The hub lives in the server's memory, so all devices of a user must reach
// the same server process. To run several backend instances, replace it with
// a shared broker (for example MongoDB change streams or Redis pub/sub).
type Hub struct {
	mu   sync.Mutex
	subs map[primitive.ObjectID]map[*subscriber]struct{}
}

// NewHub creates an empty hub.
func NewHub() *Hub {
	return &Hub{subs: map[primitive.ObjectID]map[*subscriber]struct{}{}}
}

// Subscribe opens a stream for the user. The returned channel receives the
// user's events and is closed when unsubscribe is called or the subscriber
// falls too far behind. unsubscribe is safe to call more than once.
func (h *Hub) Subscribe(userID primitive.ObjectID) (events <-chan Event, unsubscribe func()) {
	sub := &subscriber{events: make(chan Event, bufferSize)}

	h.mu.Lock()
	if h.subs[userID] == nil {
		h.subs[userID] = map[*subscriber]struct{}{}
	}
	h.subs[userID][sub] = struct{}{}
	h.mu.Unlock()

	return sub.events, func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		h.remove(userID, sub)
	}
}

// Publish sends an event to every open stream of the user. It never blocks.
func (h *Hub) Publish(userID primitive.ObjectID, event Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for sub := range h.subs[userID] {
		select {
		case sub.events <- event:
		default:
			// The client is not keeping up; drop it so it reconnects and resyncs
			h.remove(userID, sub)
		}
	}
}

// SubscriberCount returns how many streams the user has open.
func (h *Hub) SubscriberCount(userID primitive.ObjectID) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.subs[userID])
}

// remove closes and forgets a subscriber. The caller must hold h.mu.
func (h *Hub) remove(userID primitive.ObjectID, sub *subscriber) {
	subs := h.subs[userID]
	if _, ok := subs[sub]; !ok {
		return
	}
	delete(subs, sub)
	close(sub.events)
	if len(subs) == 0 {
		delete(h.subs, userID)
	}
}
