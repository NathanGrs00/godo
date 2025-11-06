package handlers

import "github.com/gin-gonic/gin"

type SettingsHandler struct{}

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{}
}

func (h *SettingsHandler) GetAll(c *gin.Context)  { c.JSON(200, "all settings") }
func (h *SettingsHandler) GetByID(c *gin.Context) { c.JSON(200, "one setting") }
func (h *SettingsHandler) Create(c *gin.Context)  { c.JSON(201, "created") }
func (h *SettingsHandler) Update(c *gin.Context)  { c.JSON(200, "updated") }
func (h *SettingsHandler) Delete(c *gin.Context)  { c.JSON(200, "deleted") }
