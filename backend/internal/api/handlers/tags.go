package handlers

import (
	"net/http"

	"github.com/NathanGrs00/godo/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagHandler struct{}

func NewTagHandler() *TagHandler {
	return &TagHandler{}
}

func (h *TagHandler) GetAll(c *gin.Context)  { c.JSON(200, "all tags") }
func (h *TagHandler) GetByID(c *gin.Context) { c.JSON(200, "one tag") }
func (h *TagHandler) Create(c *gin.Context)  { c.JSON(201, "created") }
func (h *TagHandler) Update(c *gin.Context)  { c.JSON(200, "updated") }
func (h *TagHandler) Delete(c *gin.Context)  { c.JSON(200, "deleted") }

func (h *TagHandler) DummyTags(c *gin.Context) {
	dummyTags := []models.Tag{
		{
			ID:     primitive.NewObjectID(),
			Title:  "Tag 1",
			Color:  "#FF5733",
			UserId: primitive.NewObjectID(),
		},
		{
			ID:     primitive.NewObjectID(),
			Title:  "Tag 2",
			Color:  "#33FF57",
			UserId: primitive.NewObjectID(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": dummyTags,
	})
}
