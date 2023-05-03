package controllers

import (
	"context"
	"gin-todo/configs"
	"gin-todo/models"
	"gin-todo/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var todoCollection *mongo.Collection = configs.GetCollection(configs.DB, "todos")
var validate = validator.New()

func CreateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var todo models.Todo
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, responses.TodoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&todo); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.TodoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newTodo := models.Todo{
			ID:      todo.ID,
			Text:    todo.Text,
			Checked: todo.Checked,
		}

		result, err := todoCollection.InsertOne(ctx, newTodo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TodoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.TodoResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetATodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		todoID := c.Param("todoID")
		var todo models.Todo
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(todoID)

		// find a single document
		// use the object id to find the document
		// Нужен ли еще один ID в модели?
		err := todoCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TodoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.TodoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": todo}})
	}
}

// Однако мы включили update переменную для получения обновленных полей
// и обновили коллекцию с помощью файла todoCollection.UpdateOne.
// Наконец, мы искали обновленные данные пользователя и возвращали декодированный ответ.
func EditATodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		todoId := c.Param("todoID")
		var todo models.Todo
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(todoId)

		//validate the request body
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, responses.TodoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&todo); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.TodoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"text": todo.Text, "checked": todo.Checked}
		result, err := todoCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TodoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated user details
		var updatedTodo models.Todo
		if result.MatchedCount == 1 {
			err := todoCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedTodo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.TodoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.TodoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedTodo}})
	}
}

// Функция DeleteAUserповторяет предыдущие шаги, удаляя совпадающую запись
// с помощью метода userCollection.DeleteOne.
// Мы также проверили, был ли элемент успешно удален, и вернули соответствующий ответ.
func DeleteATodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userID := c.Param("userID")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userID)

		result, err := todoCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TodoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.TodoResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.TodoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}},
		)
	}
}

func FetchAllTodos() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var todos []models.Todo
		defer cancel()

		results, err := todoCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TodoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleTodo models.Todo
			if err = results.Decode(&singleTodo); err != nil {
				c.JSON(http.StatusInternalServerError, responses.TodoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			todos = append(todos, singleTodo)
		}

		c.JSON(http.StatusOK,
			responses.TodoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": todos}},
		)
	}
}
