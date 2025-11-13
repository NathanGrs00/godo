package api

import (
	"github.com/NathanGrs00/godo/internal/api/handlers"
	_interface "github.com/NathanGrs00/godo/internal/repository/interface"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, taskRepo _interface.TaskRepository) {
	userHandlers := handlers.NewUserHandler()
	taskHandlers := handlers.NewTaskHandler(taskRepo)
	deadlinesHandler := handlers.NewDeadlineHandler()
	tagHandlers := handlers.NewTagHandler()
	settingsHandler := handlers.NewSettingsHandler()
	authHandlers := handlers.NewAuthHandler()

	r.GET("/", handlers.RootHandler)

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
		deadlines.GET("/dummies", deadlinesHandler.DummyDeadlines)
	}

	tags := r.Group("/tags")
	{
		tags.GET("/", tagHandlers.GetAll)
		tags.GET("/:id", tagHandlers.GetByID)
		tags.POST("/", tagHandlers.Create)
		tags.PUT("/:id", tagHandlers.Update)
		tags.DELETE("/:id", tagHandlers.Delete)
		tags.GET("/dummies", tagHandlers.DummyTags)
	}

	//TODO: Settings need different handlers instead of CRUD operations.
	settings := r.Group("/settings")
	{
		settings.GET("/", settingsHandler.GetAll)
		settings.GET("/:id", settingsHandler.GetByID)
		settings.POST("/", settingsHandler.Create)
		settings.PUT("/:id", settingsHandler.Update)
		settings.DELETE("/:id", settingsHandler.Delete)
	}
}
