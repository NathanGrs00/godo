package main

import (
	"github.com/NathanGrs00/godo/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api.SetupRoutes(r)

	r.Run(":8080")
}
