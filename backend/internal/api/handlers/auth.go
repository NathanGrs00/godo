package handlers

import (
	"net/http"

	"github.com/NathanGrs00/godo/internal/api/middleware"
	"github.com/NathanGrs00/godo/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler serves registration, login and the current user.
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register handles POST /auth/register and returns {"token", "user"}.
func (h *AuthHandler) Register(c *gin.Context) {
	var in service.Credentials
	if !bindJSON(c, &in) {
		return
	}
	session, err := h.service.Register(c.Request.Context(), in)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, session)
}

// Login handles POST /auth/login and returns {"token", "user"}.
func (h *AuthHandler) Login(c *gin.Context) {
	var in service.Credentials
	if !bindJSON(c, &in) {
		return
	}
	session, err := h.service.Login(c.Request.Context(), in)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, session)
}

// Me handles GET /auth/me and returns {"user"}.
func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.service.Me(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
