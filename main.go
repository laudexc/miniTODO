package main

import (
	"log"

	apihttp "miniTODO/http"
	mylib "miniTODO/myLib"
)

func main() {
	bookshelf := mylib.NewBookshelf()
	handlers := apihttp.NewHTTPHandlers(bookshelf)
	server := apihttp.NewHTTPServer(handlers)

	if err := server.StartServer(); err != nil {
		log.Fatalln(err)
	}
}
