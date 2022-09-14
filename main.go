package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{Id: "1", Title: "Clear Room", Completed: false},
	{Id: "2", Title: "Read Book", Completed: false},
	{Id: "3", Title: "Record Video", Completed: false},
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func addTodo(context *gin.Context) {
	var newTodo todo

	if err := context.BindJSON(&newTodo); err != nil {
		return
	}

	todos = append(todos, newTodo)

	context.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodo(context *gin.Context) {
	todo, shouldReturn := getTodoByIdInContext(context)
	if shouldReturn {
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func toggleTodoStatus(context *gin.Context) {
	todo, shouldReturn := getTodoByIdInContext(context)
	if shouldReturn {
		return
	}

	todo.Completed = !todo.Completed

	context.IndentedJSON(http.StatusOK, todo)
}

func removeTodo(context *gin.Context) {
	todo, shouldReturn := getTodoByIdInContext(context)
	if shouldReturn {
		return
	}

	todos = removeItemFromAList(todos, todo)

	context.IndentedJSON(http.StatusOK, todos)
}

func removeItemFromAList(list []todo, todo *todo) []todo {
	for i, t := range list {
		if t.Id == todo.Id {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func getTodoByIdInContext(context *gin.Context) (*todo, bool) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return nil, true
	}
	return todo, false
}

func getTodoById(id string) (*todo, error) {
	for i, t := range todos {
		if t.Id == id {
			return &todos[i], nil
		}
	}

	return nil, errors.New("Todo not found")
}

func main() {
	router := gin.Default()

	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)
	router.POST("/todos", addTodo)
	router.DELETE("/todos/:id", removeTodo)

	router.Run("localhost:9090")
}
