package handlers

import (
	"context"
	"net/http"

	"github.com/NathanGrs00/godo/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResourceService is the set of operations shared by the task, deadline and
// tag services. T is the model and In is the create/update input.
type ResourceService[T any, In any] interface {
	List(ctx context.Context, userID primitive.ObjectID) ([]T, error)
	Get(ctx context.Context, userID, id primitive.ObjectID) (T, error)
	Create(ctx context.Context, userID primitive.ObjectID, in In) (T, error)
	Update(ctx context.Context, userID, id primitive.ObjectID, in In) (T, error)
	Delete(ctx context.Context, userID, id primitive.ObjectID) error
}

// ResourceHandler exposes a ResourceService as REST endpoints.
// One instance serves /tasks, another /deadlines and another /tags.
type ResourceHandler[T any, In any] struct {
	service ResourceService[T, In]
	// singular and plural are the JSON keys, e.g. "task" and "tasks"
	singular string
	plural   string
}

// NewResourceHandler creates a handler that wraps responses in the given keys.
func NewResourceHandler[T any, In any](service ResourceService[T, In], singular, plural string) *ResourceHandler[T, In] {
	return &ResourceHandler[T, In]{service: service, singular: singular, plural: plural}
}

// Register adds the five CRUD routes to group.
func (h *ResourceHandler[T, In]) Register(group *gin.RouterGroup) {
	group.GET("", h.List)
	group.POST("", h.Create)
	group.GET("/:id", h.Get)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// List handles GET /<plural>.
func (h *ResourceHandler[T, In]) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{h.plural: items})
}

// Get handles GET /<plural>/:id.
func (h *ResourceHandler[T, In]) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	item, err := h.service.Get(c.Request.Context(), middleware.UserID(c), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{h.singular: item})
}

// Create handles POST /<plural>.
func (h *ResourceHandler[T, In]) Create(c *gin.Context) {
	var in In
	if !bindJSON(c, &in) {
		return
	}
	item, err := h.service.Create(c.Request.Context(), middleware.UserID(c), in)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{h.singular: item})
}

// Update handles PUT /<plural>/:id. The body must include the "version"
// the client last saw; a stale version gets 409 Conflict.
func (h *ResourceHandler[T, In]) Update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var in In
	if !bindJSON(c, &in) {
		return
	}
	item, err := h.service.Update(c.Request.Context(), middleware.UserID(c), id, in)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{h.singular: item})
}

// Delete handles DELETE /<plural>/:id.
func (h *ResourceHandler[T, In]) Delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), middleware.UserID(c), id); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
