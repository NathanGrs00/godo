// Package handlers turns HTTP requests into service calls and service
// results into JSON responses.
//
// Successful responses wrap data in a named key, e.g. {"task": {...}} or
// {"tasks": [...]}. Errors are always {"error": "message"}, plus "field"
// for validation errors.
package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/NathanGrs00/godo/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// writeError maps a service error to an HTTP status and JSON body.
func writeError(c *gin.Context, err error) {
	var validation *service.ValidationError
	switch {
	case errors.As(err, &validation):
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Error(), "field": validation.Field})
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, service.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "this item was changed on another device; reload and try again"})
	case errors.Is(err, service.ErrEmailTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		// Do not leak internal details to clients
		log.Printf("internal error on %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// bindJSON decodes the request body into dst, writing a 400 on failure.
func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return false
	}
	return true
}

// pathID parses the ":id" URL parameter, writing a 400 on failure.
func pathID(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return primitive.NilObjectID, false
	}
	return id, true
}
