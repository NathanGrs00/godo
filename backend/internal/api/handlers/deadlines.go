package handlers

import "github.com/gin-gonic/gin"

type DeadlineHandler struct{}

func NewDeadlineHandler() *DeadlineHandler {
	return &DeadlineHandler{}
}

func (h *DeadlineHandler) GetAll(c *gin.Context)  { c.JSON(200, "all deadlines") }
func (h *DeadlineHandler) GetByID(c *gin.Context) { c.JSON(200, "one deadline") }
func (h *DeadlineHandler) Create(c *gin.Context)  { c.JSON(201, "created") }
func (h *DeadlineHandler) Update(c *gin.Context)  { c.JSON(200, "updated") }
func (h *DeadlineHandler) Delete(c *gin.Context)  { c.JSON(200, "deleted") }
