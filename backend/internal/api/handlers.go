package api

import "github.com/gin-gonic/gin"

func rootHandler(c *gin.Context) {
	c.String(200, "This is the homepage.")
}
