package handlers

import (
	"net/http"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) GetAll(c *gin.Context)  { c.JSON(200, "all tasks") }
func (h *TaskHandler) GetByID(c *gin.Context) { c.JSON(200, "one task") }
func (h *TaskHandler) Create(c *gin.Context)  { c.JSON(201, "created") }
func (h *TaskHandler) Update(c *gin.Context)  { c.JSON(200, "updated") }
func (h *TaskHandler) Delete(c *gin.Context)  { c.JSON(200, "deleted") }

func (h *TaskHandler) DummyTasks(c *gin.Context) {
	dummyTasks := []models.Task{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Task 1",
			Description: "First dummy task",
			Priority:    "High",
			Completed:   false,
			UserId:      primitive.NewObjectID(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Task 2",
			Description: "Second dummy task",
			Priority:    "Medium",
			Completed:   true,
			UserId:      primitive.NewObjectID(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": dummyTasks,
	})
}
