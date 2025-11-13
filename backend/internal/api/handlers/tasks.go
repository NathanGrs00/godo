package handlers

import (
	"net/http"

	"github.com/NathanGrs00/godo/internal/models"
	_interface "github.com/NathanGrs00/godo/internal/repository/interface"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	repo _interface.TaskRepository
}

func NewTaskHandler(repo _interface.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) GetAll(c *gin.Context)  { c.JSON(200, "all tasks") }
func (h *TaskHandler) GetByID(c *gin.Context) { c.JSON(200, "one task") }

func (h *TaskHandler) Create(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = primitive.NewObjectID()

	if err := h.repo.Create(c.Request.Context(), &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    task,
	})
}

func (h *TaskHandler) Update(c *gin.Context) { c.JSON(200, "updated") }
func (h *TaskHandler) Delete(c *gin.Context) { c.JSON(200, "deleted") }

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
