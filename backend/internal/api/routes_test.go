package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/NathanGrs00/godo/internal/auth"
	"github.com/NathanGrs00/godo/internal/realtime"
	"github.com/NathanGrs00/godo/internal/repository/memory"
	"github.com/NathanGrs00/godo/internal/service"
	"github.com/gin-gonic/gin"
)

// newTestServer starts the full API on in-memory repositories.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)

	hub := realtime.NewHub()
	tokens := auth.NewTokenManager("test-secret-test-secret-test-secret", time.Hour)
	taskRepo, deadlineRepo, tagRepo := memory.NewTaskRepository(), memory.NewDeadlineRepository(), memory.NewTagRepository()
	tasks := service.NewTaskService(taskRepo, deadlineRepo, tagRepo, hub)
	deadlines := service.NewDeadlineService(deadlineRepo, tagRepo, tasks, hub)

	server := httptest.NewServer(NewRouter(Deps{
		Auth:           service.NewAuthService(memory.NewUserRepository(), tokens),
		Tasks:          tasks,
		Deadlines:      deadlines,
		Tags:           service.NewTagService(tagRepo, tasks, deadlines, hub),
		Tokens:         tokens,
		Events:         hub,
		AllowedOrigins: []string{"http://localhost:5173"},
		Heartbeat:      50 * time.Millisecond,
	}))
	t.Cleanup(server.Close)
	return server
}

// client is a logged-in device.
type client struct {
	t      *testing.T
	server *httptest.Server
	token  string
}

// do sends a JSON request and decodes the JSON response into out (if not nil).
func (c *client) do(method, path string, body any, out any) int {
	c.t.Helper()
	var reader io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	}
	req, _ := http.NewRequest(method, c.server.URL+path, reader)
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.t.Fatal(err)
	}
	defer res.Body.Close()
	if out != nil && res.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			c.t.Fatalf("%s %s: decode: %v", method, path, err)
		}
	}
	return res.StatusCode
}

// register creates an account and returns a client logged in to it.
func register(t *testing.T, server *httptest.Server, email string) *client {
	c := &client{t: t, server: server}
	var session struct{ Token string }
	if code := c.do("POST", "/auth/register", map[string]string{"email": email, "password": "password123"}, &session); code != http.StatusCreated {
		t.Fatalf("register status = %d", code)
	}
	c.token = session.Token
	return c
}

// login returns a second client (device) for an existing account.
func login(t *testing.T, server *httptest.Server, email string) *client {
	c := &client{t: t, server: server}
	var session struct{ Token string }
	if code := c.do("POST", "/auth/login", map[string]string{"email": email, "password": "password123"}, &session); code != http.StatusOK {
		t.Fatalf("login status = %d", code)
	}
	c.token = session.Token
	return c
}

type taskBody struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Priority  string   `json:"priority"`
	Completed bool     `json:"completed"`
	TagIDs    []string `json:"tag_ids"`
	Version   int64    `json:"version"`
}

func TestHealthAndAuth(t *testing.T) {
	server := newTestServer(t)
	anon := &client{t: t, server: server}

	var health map[string]string
	if code := anon.do("GET", "/health", nil, &health); code != 200 || health["status"] != "ok" {
		t.Fatalf("health = %d %v", code, health)
	}

	if code := anon.do("GET", "/tasks", nil, nil); code != http.StatusUnauthorized {
		t.Fatalf("no token: status %d, want 401", code)
	}
	bad := &client{t: t, server: server, token: "nope"}
	if code := bad.do("GET", "/tasks", nil, nil); code != http.StatusUnauthorized {
		t.Fatalf("bad token: status %d, want 401", code)
	}

	alice := register(t, server, "alice@example.com")
	var me struct{ User struct{ Email string } }
	if code := alice.do("GET", "/auth/me", nil, &me); code != 200 || me.User.Email != "alice@example.com" {
		t.Fatalf("me = %d %+v", code, me)
	}

	var errBody map[string]string
	if code := anon.do("POST", "/auth/register", map[string]string{"email": "alice@example.com", "password": "password123"}, &errBody); code != http.StatusConflict {
		t.Fatalf("duplicate register status = %d", code)
	}
	if code := anon.do("POST", "/auth/login", map[string]string{"email": "alice@example.com", "password": "wrong-password"}, &errBody); code != http.StatusUnauthorized {
		t.Fatalf("bad login status = %d", code)
	}
}

func TestTaskCRUD(t *testing.T) {
	server := newTestServer(t)
	alice := register(t, server, "alice@example.com")

	var tagRes struct{ Tag struct{ ID string } }
	if code := alice.do("POST", "/tags", map[string]string{"title": "work", "color": "#8CFBDE"}, &tagRes); code != http.StatusCreated {
		t.Fatalf("create tag status = %d", code)
	}

	var created struct{ Task taskBody }
	code := alice.do("POST", "/tasks", map[string]any{"title": "Ship it", "priority": "high", "tag_ids": []string{tagRes.Tag.ID}}, &created)
	if code != http.StatusCreated || created.Task.Version != 1 || len(created.Task.TagIDs) != 1 {
		t.Fatalf("create = %d %+v", code, created)
	}
	id := created.Task.ID

	var list struct{ Tasks []taskBody }
	if alice.do("GET", "/tasks", nil, &list); len(list.Tasks) != 1 {
		t.Fatalf("list = %+v", list)
	}

	var updated struct{ Task taskBody }
	code = alice.do("PUT", "/tasks/"+id, map[string]any{"title": "Ship it", "priority": "high", "completed": true, "version": 1}, &updated)
	if code != 200 || !updated.Task.Completed || updated.Task.Version != 2 || len(updated.Task.TagIDs) != 0 {
		t.Fatalf("update = %d %+v", code, updated)
	}

	// A device still holding version 1 is told to reload
	var conflict map[string]string
	if code := alice.do("PUT", "/tasks/"+id, map[string]any{"title": "Old", "priority": "low", "version": 1}, &conflict); code != http.StatusConflict {
		t.Fatalf("stale update status = %d", code)
	}

	var validation map[string]string
	if code := alice.do("POST", "/tasks", map[string]any{"title": "", "priority": "high"}, &validation); code != http.StatusBadRequest || validation["field"] != "title" {
		t.Fatalf("validation = %d %v", code, validation)
	}
	if code := alice.do("POST", "/tasks", "not an object", &validation); code != http.StatusBadRequest {
		t.Fatalf("bad body status = %d", code)
	}
	if code := alice.do("GET", "/tasks/not-an-id", nil, &validation); code != http.StatusBadRequest {
		t.Fatalf("bad id status = %d", code)
	}

	// Other users cannot see or touch it
	bob := register(t, server, "bob@example.com")
	if code := bob.do("GET", "/tasks/"+id, nil, &validation); code != http.StatusNotFound {
		t.Fatalf("bob get status = %d", code)
	}
	if code := bob.do("DELETE", "/tasks/"+id, nil, nil); code != http.StatusNotFound {
		t.Fatalf("bob delete status = %d", code)
	}

	if code := alice.do("DELETE", "/tasks/"+id, nil, nil); code != http.StatusNoContent {
		t.Fatalf("delete status = %d", code)
	}
	if code := alice.do("GET", "/tasks/"+id, nil, &validation); code != http.StatusNotFound {
		t.Fatalf("get after delete status = %d", code)
	}
}

func TestDeadlinesAndTagsRoutes(t *testing.T) {
	server := newTestServer(t)
	alice := register(t, server, "alice@example.com")

	var d struct {
		Deadline struct {
			ID     string
			Passed bool
			Date   time.Time
		}
	}
	date := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	if code := alice.do("POST", "/deadlines", map[string]any{"title": "Exam", "date": date}, &d); code != http.StatusCreated || d.Deadline.Passed || !d.Deadline.Date.Equal(date) {
		t.Fatalf("create deadline = %d %+v", code, d)
	}
	var list struct{ Deadlines []any }
	if alice.do("GET", "/deadlines", nil, &list); len(list.Deadlines) != 1 {
		t.Fatalf("deadlines = %+v", list)
	}
	var tags struct{ Tags []any }
	if alice.do("GET", "/tags", nil, &tags); tags.Tags == nil {
		t.Fatal("empty tag list should be [] not null")
	}
}

func TestCORS(t *testing.T) {
	server := newTestServer(t)
	req, _ := http.NewRequest("OPTIONS", server.URL+"/tasks", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "PUT")
	req.Header.Set("Access-Control-Request-Headers", "Authorization")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.Header.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("preflight headers = %v", res.Header)
	}
}

// TestSyncBetweenDevices is the core sync scenario: a change made on one
// device arrives on another device of the same user, and not on anyone else's.
func TestSyncBetweenDevices(t *testing.T) {
	server := newTestServer(t)
	laptop := register(t, server, "alice@example.com")
	phone := login(t, server, "alice@example.com")
	bob := register(t, server, "bob@example.com")

	phoneEvents := openStream(t, server, phone.token)
	bobEvents := openStream(t, server, bob.token)

	var created struct{ Task taskBody }
	laptop.do("POST", "/tasks", map[string]any{"title": "Buy milk", "priority": "low"}, &created)

	event := phoneEvents.next(t)
	if event.Type != realtime.TaskCreated || event.ID != created.Task.ID {
		t.Fatalf("phone got %+v", event)
	}
	var data taskBody
	json.Unmarshal(event.Data, &data)
	if data.Title != "Buy milk" {
		t.Fatalf("event data = %+v", data)
	}

	laptop.do("DELETE", "/tasks/"+created.Task.ID, nil, nil)
	if event := phoneEvents.next(t); event.Type != realtime.TaskDeleted {
		t.Fatalf("phone got %+v", event)
	}

	bobEvents.expectNone(t, 150*time.Millisecond)
}

// stream reads Server-Sent Events from /events.
type stream struct {
	events chan streamEvent
}

type streamEvent struct {
	Type string          `json:"type"`
	ID   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}

// openStream connects to /events and waits until the server says it is ready.
func openStream(t *testing.T, server *httptest.Server, token string) *stream {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL+"/events", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("content type = %q", ct)
	}

	s := &stream{events: make(chan streamEvent, 16)}
	ready := make(chan struct{})
	go func() {
		defer res.Body.Close()
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "event: ready" {
				close(ready)
				continue
			}
			if payload, ok := strings.CutPrefix(line, "data: "); ok && payload != "{}" {
				var e streamEvent
				if json.Unmarshal([]byte(payload), &e) == nil {
					s.events <- e
				}
			}
		}
	}()

	select {
	case <-ready:
	case <-time.After(2 * time.Second):
		t.Fatal("stream never became ready")
	}
	return s
}

func (s *stream) next(t *testing.T) streamEvent {
	t.Helper()
	select {
	case e := <-s.events:
		return e
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for event")
		return streamEvent{}
	}
}

func (s *stream) expectNone(t *testing.T, wait time.Duration) {
	t.Helper()
	select {
	case e := <-s.events:
		t.Fatalf("unexpected event %+v", e)
	case <-time.After(wait):
	}
}
