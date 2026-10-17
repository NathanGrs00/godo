// Package api wires handlers, middleware and routes into a Gin engine.
package api

import (
	"net/http"
	"time"

	"github.com/NathanGrs00/godo/internal/api/handlers"
	"github.com/NathanGrs00/godo/internal/api/middleware"
	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Deps is everything the router needs.
type Deps struct {
	Auth      *service.AuthService
	Tasks     *service.TaskService
	Deadlines *service.DeadlineService
	Tags      *service.TagService
	Tokens    middleware.TokenVerifier
	Events    handlers.Subscriber
	// AllowedOrigins are the frontend URLs allowed by CORS.
	AllowedOrigins []string
	// Heartbeat is how often the event stream sends a keep-alive (default 25s).
	Heartbeat time.Duration
}

// NewRouter builds the HTTP API.
//
//	GET    /health                 liveness check (public)
//	POST   /auth/register          create account (public)
//	POST   /auth/login             log in (public)
//	GET    /auth/me                current user
//	GET    /events                 live change stream (Server-Sent Events)
//	CRUD   /tasks, /deadlines, /tags
//
// Every route except the public ones needs "Authorization: Bearer <token>".
func NewRouter(deps Deps) *gin.Engine {
	if deps.Heartbeat == 0 {
		deps.Heartbeat = 25 * time.Second
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: deps.AllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authHandler := handlers.NewAuthHandler(deps.Auth)
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// Everything below requires a logged-in user
	private := r.Group("", middleware.RequireAuth(deps.Tokens))
	private.GET("/auth/me", authHandler.Me)
	private.GET("/events", handlers.NewEventsHandler(deps.Events, deps.Heartbeat).Stream)

	handlers.NewResourceHandler[*models.Task, service.TaskInput](deps.Tasks, "task", "tasks").
		Register(private.Group("/tasks"))
	handlers.NewResourceHandler[*models.Deadline, service.DeadlineInput](deps.Deadlines, "deadline", "deadlines").
		Register(private.Group("/deadlines"))
	handlers.NewResourceHandler[*models.Tag, service.TagInput](deps.Tags, "tag", "tags").
		Register(private.Group("/tags"))

	return r
}
