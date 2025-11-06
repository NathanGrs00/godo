package handlers

import "github.com/gin-gonic/gin"

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) GetAll(c *gin.Context)  { c.JSON(200, "all tasks") }
func (h *TaskHandler) GetByID(c *gin.Context) { c.JSON(200, "one task") }
func (h *TaskHandler) Create(c *gin.Context)  { c.JSON(201, "created") }
func (h *TaskHandler) Update(c *gin.Context)  { c.JSON(200, "updated") }
func (h *TaskHandler) Delete(c *gin.Context)  { c.JSON(200, "deleted") }
