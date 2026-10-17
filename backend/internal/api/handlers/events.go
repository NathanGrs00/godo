package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NathanGrs00/godo/internal/api/middleware"
	"github.com/NathanGrs00/godo/internal/realtime"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Subscriber opens a user's event stream. realtime.Hub implements it.
type Subscriber interface {
	Subscribe(userID primitive.ObjectID) (<-chan realtime.Event, func())
}

// EventsHandler streams change events to a device using Server-Sent Events.
//
// Stream format (one message per change):
//
//	event: ready            sent once the stream is open
//	data: {"type":"task.updated","id":"...","data":{...}}
//	: ping                  keep-alive comment every heartbeat interval
type EventsHandler struct {
	hub       Subscriber
	heartbeat time.Duration
}

// NewEventsHandler creates an EventsHandler that sends a keep-alive every heartbeat.
func NewEventsHandler(hub Subscriber, heartbeat time.Duration) *EventsHandler {
	return &EventsHandler{hub: hub, heartbeat: heartbeat}
}

// Stream handles GET /events. It stays open until the client disconnects.
func (h *EventsHandler) Stream(c *gin.Context) {
	events, unsubscribe := h.hub.Subscribe(middleware.UserID(c))
	defer unsubscribe()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// Stop reverse proxies such as nginx from buffering the stream
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	// Tell the client the subscription is live, so it can reload any data
	// it might have missed while it was disconnected.
	fmt.Fprint(c.Writer, "event: ready\ndata: {}\n\n")
	c.Writer.Flush()

	ticker := time.NewTicker(h.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			fmt.Fprint(c.Writer, ": ping\n\n")
			c.Writer.Flush()
		case event, open := <-events:
			if !open {
				// The hub dropped us (we fell behind); the client will reconnect
				return
			}
			payload, err := json.Marshal(event)
			if err != nil {
				continue
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", payload)
			c.Writer.Flush()
		}
	}
}
