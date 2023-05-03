package main

import (
	"gin-todo/configs"
	"gin-todo/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// run database connection
	configs.ConnectDB()

	// run routes
	routes.TodoRoute(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "Hello Hello from Gin-gonic & mongoDB",
		})
	})

	router.Run("localhost:8080")
}
