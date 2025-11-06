package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", rootHandler)
	r.GET("/tasks", tasksHandler)
	r.GET("/users", usersHandler)
	r.POST("/login", loginHandler)
	r.POST("/register", registerHandler)
	r.GET("/deadlines", deadlinesHandler)
	r.GET("/settings", settingsHandler)
}
