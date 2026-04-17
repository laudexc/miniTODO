package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandler *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandler,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/books").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateBook)
	router.Path("/books").Methods("GET").Queries("title", "{title}").HandlerFunc(s.httpHandlers.HandleGetByTitle)
	router.Path("/books").Methods("GET").Queries("author", "{author}").HandlerFunc(s.httpHandlers.HandleGetByAuthor)
	router.Path("/books").Methods("GET").Queries("completed", "true").HandlerFunc(s.httpHandlers.HandleGetAllReadedBooks)
	router.Path("/books").Methods("GET").Queries("completed", "false").HandlerFunc(s.httpHandlers.HandleGetAllUnreadedBooks)
	router.Path("/books").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetAllBooks)
	router.Path("/books/{title}").Methods("PATCH").HandlerFunc(s.httpHandlers.HandleReadBook)
	router.Path("/books/{title}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteBook)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}
