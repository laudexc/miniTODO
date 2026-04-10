package main

import (
	"log"

	apihttp "restAPI/http"
	"restAPI/todo"
)

func main() {
	todoList := todo.NewList()
	handlers := apihttp.NewHTTPHandlers(todoList)
	server := apihttp.NewHTTPServer(handlers)

	if err := server.StartServer(); err != nil {
		log.Fatalln(err)
	}
}
