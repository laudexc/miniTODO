package main

import (
	"log"

	apihttp "httpLibrary/http"
	"httpLibrary/todo"
)

func main() {
	todoList := todo.NewList()
	handlers := apihttp.NewHTTPHandlers(todoList)
	server := apihttp.NewHTTPServer(handlers)

	if err := server.StartServer(); err != nil {
		log.Fatalln(err)
	}
}
