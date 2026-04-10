package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"restAPI/todo"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
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
func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(todoTask, "", "    ")
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
	 pattern: /tasks/{title}
	 method: GET
	 info: pattern

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented found task

	 failed:
		- status code: 400, 404, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"] // вернет мапу с ключем - параметром, значением - что там написал клиент;
	// ок - написал ли что то или нет, можно не проверять, ибо если этот хендлер уже вызвался - значит что то было передано

	task, err := h.todoList.GetTask(title)
	if err != nil {
		// errDTO := ErrorDTO {
		// 	Message: err.Error(),
		// 	Time: time.Now(),
		// }
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
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
func (h HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListTasks()
	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

/*
	 pattern: /tasks?completed=false
	 method: GET
	 info: query params

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented found tasks

	 failed:
		- status code: 400, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleGetAllUncomplitedTAsks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks := h.todoList.ListUncompletedTasks()
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
	 pattern: /tasks/{title}
	 method: PATCH
	 info: pattern + JSON in request body

	 succeed:
		- status code: 200 OK
		- response body: JSON respresented changed tasks

	 failed:
		- status code: 400, 409, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO completeTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]
	var task todo.Task
	var err error

	if completeDTO.Complete {
		task, err = h.todoList.CompleteTask(title)
		if err != nil {
			errDTO := NewErrorDTO(err)
			if errors.Is(err, todo.ErrTaskNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
			} else {
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}

			return
		}
	} else {
		task, err = h.todoList.UncompleteTask(title)
		if err != nil {
			errDTO := NewErrorDTO(err)
			if errors.Is(err, todo.ErrTaskNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
			} else {
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}

			return
		}
	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

/*
	 pattern: /tasks/{title}
	 method: DELETE
	 info: pattern

	 succeed:
		- status code: 204 No Content
		- response body: -

	 failed:
		- status code: 400, 404, 500...
		- response body: JSON with error + time
*/
func (h HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	if err := h.todoList.DeleteTask(title); err != nil {
		errDTO := NewErrorDTO(err)
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
