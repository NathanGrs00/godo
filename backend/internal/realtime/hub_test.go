package realtime

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPublishReachesOnlyThatUsersDevices(t *testing.T) {
	hub := NewHub()
	alice, bob := primitive.NewObjectID(), primitive.NewObjectID()
	laptop, closeLaptop := hub.Subscribe(alice)
	phone, closePhone := hub.Subscribe(alice)
	bobs, closeBob := hub.Subscribe(bob)
	defer closeLaptop()
	defer closePhone()
	defer closeBob()

	hub.Publish(alice, Event{Type: TaskCreated, ID: "1"})

	for name, ch := range map[string]<-chan Event{"laptop": laptop, "phone": phone} {
		select {
		case e := <-ch:
			if e.ID != "1" {
				t.Errorf("%s got %+v", name, e)
			}
		default:
			t.Errorf("%s got nothing", name)
		}
	}
	select {
	case e := <-bobs:
		t.Fatalf("bob received alice's event %+v", e)
	default:
	}
}

func TestUnsubscribe(t *testing.T) {
	hub := NewHub()
	user := primitive.NewObjectID()
	events, unsubscribe := hub.Subscribe(user)
	if hub.SubscriberCount(user) != 1 {
		t.Fatal("subscriber not counted")
	}
	unsubscribe()
	unsubscribe() // safe twice
	if hub.SubscriberCount(user) != 0 {
		t.Fatal("subscriber not removed")
	}
	if _, open := <-events; open {
		t.Fatal("channel should be closed")
	}
	hub.Publish(user, Event{Type: TaskDeleted}) // no subscribers: must not panic or block
}

func TestSlowSubscriberIsDropped(t *testing.T) {
	hub := NewHub()
	user := primitive.NewObjectID()
	events, unsubscribe := hub.Subscribe(user)
	defer unsubscribe()

	for i := 0; i <= bufferSize; i++ {
		hub.Publish(user, Event{Type: TaskUpdated})
	}
	if hub.SubscriberCount(user) != 0 {
		t.Fatal("slow subscriber should have been dropped")
	}
	count := 0
	for range events {
		count++
	}
	if count != bufferSize {
		t.Fatalf("drained %d buffered events, want %d", count, bufferSize)
	}
}
