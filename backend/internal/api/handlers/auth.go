package handlers

import "github.com/gin-gonic/gin"

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "This is the login endpoint."})
}

func (h *AuthHandler) Register(c *gin.Context) {
	c.JSON(200, gin.H{"message": "This is the register endpoint."})
}
