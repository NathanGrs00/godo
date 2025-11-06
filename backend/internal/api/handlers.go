package api

import "github.com/gin-gonic/gin"

func rootHandler(c *gin.Context) {
	c.String(200, "This is the homepage.")
}

func tasksHandler(c *gin.Context) {
	c.String(200, "This is the tasks endpoint.")
}

func usersHandler(c *gin.Context) {
	c.String(200, "This is the users endpoint.")
}

func loginHandler(c *gin.Context) {
	c.String(200, "This is the login endpoint.")
}

func registerHandler(c *gin.Context) {
	c.String(200, "This is the register endpoint.")
}

func deadlinesHandler(c *gin.Context) {
	c.String(200, "This is the deadlines endpoint.")
}

func settingsHandler(c *gin.Context) {
	c.String(200, "This is the settings endpoint.")
}
