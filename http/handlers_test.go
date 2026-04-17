package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mylib "miniTODO/myLib"

	"github.com/gorilla/mux"
)

func newTestRouter() *mux.Router {
	bookshelf := mylib.NewBookshelf()
	handlers := NewHTTPHandlers(bookshelf)

	router := mux.NewRouter()
	router.Path("/books").Methods("POST").HandlerFunc(handlers.HandleCreateBook)
	router.Path("/books").Methods("GET").Queries("title", "{title}").HandlerFunc(handlers.HandleGetByTitle)
	router.Path("/books").Methods("GET").Queries("author", "{author}").HandlerFunc(handlers.HandleGetByAuthor)
	router.Path("/books").Methods("GET").Queries("completed", "true").HandlerFunc(handlers.HandleGetAllReadedBooks)
	router.Path("/books").Methods("GET").Queries("completed", "false").HandlerFunc(handlers.HandleGetAllUnreadedBooks)
	router.Path("/books").Methods("GET").HandlerFunc(handlers.HandleGetAllBooks)
	router.Path("/books/{title}").Methods("PATCH").HandlerFunc(handlers.HandleReadBook)
	router.Path("/books/{title}").Methods("DELETE").HandlerFunc(handlers.HandleDeleteBook)

	return router
}

func doRequest(t *testing.T, router http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func createBook(t *testing.T, router http.Handler, title, author string, pages int) *httptest.ResponseRecorder {
	t.Helper()

	payload := map[string]any{
		"Title":      title,
		"Author":     author,
		"NumOfPages": pages,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal create payload: %v", err)
	}
	return doRequest(t, router, http.MethodPost, "/books", b)
}

func TestCreateBook_Success(t *testing.T) {
	router := newTestRouter()
	rr := createBook(t, router, "1984", "George Orwell", 328)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusCreated, rr.Code, rr.Body.String())
	}
}

func TestCreateBook_Duplicate(t *testing.T) {
	router := newTestRouter()
	_ = createBook(t, router, "1984", "George Orwell", 328)
	rr := createBook(t, router, "1984", "George Orwell", 328)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusConflict, rr.Code, rr.Body.String())
	}
}

func TestGetByTitle_Found(t *testing.T) {
	router := newTestRouter()
	_ = createBook(t, router, "1984", "George Orwell", 328)

	rr := doRequest(t, router, http.MethodGet, "/books?title=1984", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}
}

func TestGetByTitle_NotFound(t *testing.T) {
	router := newTestRouter()
	rr := doRequest(t, router, http.MethodGet, "/books?title=missing", nil)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestGetByAuthor_Filter(t *testing.T) {
	router := newTestRouter()
	_ = createBook(t, router, "1984", "George Orwell", 328)
	_ = createBook(t, router, "Animal Farm", "George Orwell", 112)
	_ = createBook(t, router, "Dune", "Frank Herbert", 412)

	rr := doRequest(t, router, http.MethodGet, "/books?author=George+Orwell", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var got map[string]mylib.Book
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 books by author, got %d", len(got))
	}
}

func TestGetAllBooks(t *testing.T) {
	router := newTestRouter()
	_ = createBook(t, router, "1984", "George Orwell", 328)
	_ = createBook(t, router, "Dune", "Frank Herbert", 412)

	rr := doRequest(t, router, http.MethodGet, "/books", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var got map[string]mylib.Book
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 books, got %d", len(got))
	}
}

func TestReadAndCompletedFilters(t *testing.T) {
	router := newTestRouter()
	_ = createBook(t, router, "1984", "George Orwell", 328)
	_ = createBook(t, router, "Dune", "Frank Herbert", 412)

	patchBody := []byte(`{"complete": true}`)
	patchRR := doRequest(t, router, http.MethodPatch, "/books/1984", patchBody)
	if patchRR.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, patchRR.Code, patchRR.Body.String())
	}

	readRR := doRequest(t, router, http.MethodGet, "/books?completed=true", nil)
	if readRR.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, readRR.Code, readRR.Body.String())
	}

	var readBooks map[string]mylib.Book
	if err := json.Unmarshal(readRR.Body.Bytes(), &readBooks); err != nil {
		t.Fatalf("unmarshal read books: %v", err)
	}
	if len(readBooks) != 1 {
		t.Fatalf("expected 1 read book, got %d", len(readBooks))
	}

	unreadRR := doRequest(t, router, http.MethodGet, "/books?completed=false", nil)
	if unreadRR.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, unreadRR.Code, unreadRR.Body.String())
	}

	var unreadBooks map[string]mylib.Book
	if err := json.Unmarshal(unreadRR.Body.Bytes(), &unreadBooks); err != nil {
		t.Fatalf("unmarshal unread books: %v", err)
	}
	if len(unreadBooks) != 1 {
		t.Fatalf("expected 1 unread book, got %d", len(unreadBooks))
	}
}

func TestReadBook_NotFound(t *testing.T) {
	router := newTestRouter()
	rr := doRequest(t, router, http.MethodPatch, "/books/missing", []byte(`{"completed": true}`))

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestDeleteBook_SuccessAndNotFound(t *testing.T) {
	router := newTestRouter()
	_ = createBook(t, router, "1984", "George Orwell", 328)

	delRR := doRequest(t, router, http.MethodDelete, "/books/1984", nil)
	if delRR.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNoContent, delRR.Code, delRR.Body.String())
	}

	delAgainRR := doRequest(t, router, http.MethodDelete, "/books/1984", nil)
	if delAgainRR.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, delAgainRR.Code, delAgainRR.Body.String())
	}
}
