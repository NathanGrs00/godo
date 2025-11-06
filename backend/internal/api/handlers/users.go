package handlers

import "github.com/gin-gonic/gin"

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetAll(c *gin.Context)  { c.JSON(200, "all users") }
func (h *UserHandler) GetByID(c *gin.Context) { c.JSON(200, "one user") }
func (h *UserHandler) Create(c *gin.Context)  { c.JSON(201, "created") }
func (h *UserHandler) Update(c *gin.Context)  { c.JSON(200, "updated") }
func (h *UserHandler) Delete(c *gin.Context)  { c.JSON(200, "deleted") }
