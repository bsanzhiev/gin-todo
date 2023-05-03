package routes

import "github.com/gin-gonic/gin"
import "gin-todo/controllers"

func TodoRoute(router *gin.Engine)  {
    //All routes related to users comes here
    router.POST("/todo", controllers.CreateTodo())
    router.GET("/todo/:todoID", controllers.GetATodo())
    router.PUT("/todo/:todoID", controllers.EditATodo())
    router.DELETE("/todo/:todoID", controllers.DeleteATodo())
    router.GET("/todos", controllers.FetchAllTodos())
}