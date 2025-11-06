package api

import (
	"github.com/NathanGrs00/godo/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	userHandlers := handlers.NewUserHandler()
	taskHandlers := handlers.NewTaskHandler()
	deadlinesHandler := handlers.NewDeadlineHandler()
	settingsHandler := handlers.NewSettingsHandler()
	authHandlers := handlers.NewAuthHandler()

	r.GET("/", rootHandler)

	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandlers.Login)
		auth.POST("/register", authHandlers.Register)
	}

	users := r.Group("/users")
	{
		users.GET("/", userHandlers.GetAll)
		users.GET("/:id", userHandlers.GetByID)
		users.POST("/", userHandlers.Create)
		users.PUT("/:id", userHandlers.Update)
		users.DELETE("/:id", userHandlers.Delete)
	}

	tasks := r.Group("/tasks")
	{
		tasks.GET("/", taskHandlers.GetAll)
		tasks.GET("/:id", taskHandlers.GetByID)
		tasks.POST("/", taskHandlers.Create)
		tasks.PUT("/:id", taskHandlers.Update)
		tasks.DELETE("/:id", taskHandlers.Delete)
	}

	deadlines := r.Group("/deadlines")
	{
		deadlines.GET("/", deadlinesHandler.GetAll)
		deadlines.GET("/:id", deadlinesHandler.GetByID)
		deadlines.POST("/", deadlinesHandler.Create)
		deadlines.PUT("/:id", deadlinesHandler.Update)
		deadlines.DELETE("/:id", deadlinesHandler.Delete)
	}

	settings := r.Group("/settings")
	{
		settings.GET("/", settingsHandler.GetAll)
		settings.GET("/:id", settingsHandler.GetByID)
		settings.POST("/", settingsHandler.Create)
		settings.PUT("/:id", settingsHandler.Update)
		settings.DELETE("/:id", settingsHandler.Delete)
	}
}
