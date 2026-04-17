package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"miniTODO/myLib"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	bookshelf *myLib.Bookshelf
}

func NewHTTPHandlers(BookShlf *myLib.Bookshelf) *HTTPHandlers {
	return &HTTPHandlers{
		bookshelf: BookShlf,
	}
}

/*
pattern: /tasks
method:  POST
info:    JSON in HTTP request body

succeed:
  - status code: 201 Created
  - response body: JSON respresent created task

failed:
  - status code: 400, 409, 500...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	var bookDTO BookDTO
	if err := json.NewDecoder(r.Body).Decode(&bookDTO); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := bookDTO.ValidateForCreate(); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	book := myLib.NewBook(bookDTO.Title, bookDTO.Author, bookDTO.NumOfPages)
	if err := h.bookshelf.AddBook(book); err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, myLib.ErrBookAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
	 pattern: /tasks
	 method: GET
	 info: -

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented found tasks

	 failed:
		- status code: 400, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	books := h.bookshelf.ListBooks()
	b, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

/*
	 pattern: /tasks?title=string
	 method: GET
	 info: query params

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented found task

	 failed:
		- status code: 400, 404, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleGetByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		errDTO := NewErrorDTO(errors.New("query param 'title' is required"))
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	book, err := h.bookshelf.GetBook(title)
	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, myLib.ErrBookNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

/*
	 pattern: /tasks?author=string
	 method: GET
	 info: query params

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented found task

	 failed:
		- status code: 400, 404, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleGetByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if author == "" {
		errDTO := NewErrorDTO(errors.New("query param 'author' is required"))
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	books := h.bookshelf.ListByAuthorBook(author)

	b, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

/*
	 pattern: /books?completed=true
	 method: GET
	 info: query params

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented found tasks

	 failed:
		- status code: 400, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleGetAllReadedBooks(w http.ResponseWriter, r *http.Request) {
	completedTasks := h.bookshelf.ListCompletedBooks()
	b, err := json.MarshalIndent(completedTasks, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

/*
	 pattern: /books?completed=false
	 method: GET
	 info: query params

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented found tasks

	 failed:
		- status code: 400, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleGetAllUnreadedBooks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks := h.bookshelf.ListUncompletedBooks()
	b, err := json.MarshalIndent(uncompletedTasks, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}

}

/*
	 pattern: /books/{title}
	 method: PATCH
	 info: pattern + JSON in request body

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented changed tasks

	 failed:
		- status code: 400, 409, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleReadBook(w http.ResponseWriter, r *http.Request) {
	var completeDTO completedBookDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	shouldComplete, err := completeDTO.CompletionValue()
	if err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]
	var (
		changedBook myLib.Book
		updateErr   error
	)

	if shouldComplete {
		changedBook, updateErr = h.bookshelf.ReadBook(title)
	} else {
		changedBook, updateErr = h.bookshelf.UnreadBook(title)
	}

	if updateErr != nil {
		errDTO := NewErrorDTO(updateErr)

		if errors.Is(updateErr, myLib.ErrBookNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(changedBook, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

/*
	 pattern: /books/{title}
	 method: DELETE
	 info: pattern

	 succeed:
		- status code: 204 No Content
		- response body: -

	 failed:
		- status code: 400, 404, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	if err := h.bookshelf.DeleteBook(title); err != nil {
		errDTO := NewErrorDTO(err)
		if errors.Is(err, myLib.ErrBookNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
