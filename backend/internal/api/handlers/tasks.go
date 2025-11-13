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

func (h *TaskHandler) GetAll(c *gin.Context) {
	tasks, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "All tasks fetched successfully",
		"tasks":   tasks,
	})
}

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
func (h *TaskHandler) Delete(c *gin.Context) {
	idParam := c.Param("id") // get id from URL
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	err = h.repo.Delete(c.Request.Context(), objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task for id " + idParam})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})

}
