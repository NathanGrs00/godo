package handlers

import (
	"net/http"
	"time"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeadlineHandler struct{}

func NewDeadlineHandler() *DeadlineHandler {
	return &DeadlineHandler{}
}

func (h *DeadlineHandler) GetAll(c *gin.Context)  { c.JSON(200, "all deadlines") }
func (h *DeadlineHandler) GetByID(c *gin.Context) { c.JSON(200, "one deadline") }
func (h *DeadlineHandler) Create(c *gin.Context)  { c.JSON(201, "created") }
func (h *DeadlineHandler) Update(c *gin.Context)  { c.JSON(200, "updated") }
func (h *DeadlineHandler) Delete(c *gin.Context)  { c.JSON(200, "deleted") }

func (h *TaskHandler) DummyDeadlines(c *gin.Context) {
	dummyDeadlines := []models.Deadline{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Deadline 1",
			Description: "First dummy deadline",
			Date:        time.Now().Add(48 * time.Hour),
			Passed:      false,
			UserId:      primitive.NewObjectID(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Deadline 2",
			Description: "First dummy deadline",
			Date:        time.Now().Add(48 * time.Hour),
			Passed:      false,
			UserId:      primitive.NewObjectID(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"deadlines": dummyDeadlines,
	})
}
