package todos

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"todo-list-api/internal/auth"
	todosDB "todo-list-api/internal/db/todos"
	"todo-list-api/internal/models/todos"
	"todo-list-api/internal/models/user"
)

type PageData struct {
	Todos        []todos.Todo
	CurrentPage  int
	TotalPages   int
	PreviousPage int
	NextPage     int
}

func HandlerTodos(w http.ResponseWriter, r *http.Request) {
	if userInfo, ok := auth.HandleLoginWithToken(w, r); ok != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	} else if r.Method == http.MethodGet {
		renderTodosPage(w, r, userInfo)
	} else if r.Method == http.MethodPost {
		handlerTodosPost(w, r, userInfo)
	} else if r.Method == http.MethodPut {
		handlerTodosPut(w, r, userInfo)
	} else if r.Method == http.MethodDelete {
		handlerTodosDelete(w, r, userInfo)
	} else if r.Method == http.MethodPatch {
		handlerTodosPatch(w, r, userInfo)
	}
}

func renderTodosPage(w http.ResponseWriter, r *http.Request, userInfo user.Info) {
	tmpl, err := template.ParseFiles("web/templates/todos.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	limit := 5 // default limit
	page := 1  // default value
	totalTodos, err := todosDB.GetCountTodos(userInfo)
	totalPages := (totalTodos + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)
		if err == nil && parsedPage > 0 && parsedPage <= totalPages {
			page = parsedPage
		}
	}

	offset := (page - 1) * limit
	listTodos, err := todosDB.GetTodosWithPagination(userInfo, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Todos:        listTodos,
		CurrentPage:  page,
		TotalPages:   totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
	}

	if err := tmpl.Execute(w, pageData); err != nil {
		log.Println(err)
	}
}

func handlerTodosPost(w http.ResponseWriter, r *http.Request, userInfo user.Info) {
	var todoTitleDesc struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	err := json.NewDecoder(r.Body).Decode(&todoTitleDesc)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusForbidden)
		return
	}

	todo := todos.Todo{
		ID:          "0",
		UserId:      strconv.Itoa(userInfo.ID),
		Title:       todoTitleDesc.Title,
		Description: todoTitleDesc.Description,
		IsDone:      false,
	}

	todo, err = todosDB.AddTodo(todo)
	if err != nil {
		http.Error(w, "Error add todo.", http.StatusForbidden)
		return
	}

	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Invalid response data.", http.StatusForbidden)
		return
	}
}

func handlerTodosPut(w http.ResponseWriter, r *http.Request, userInfo user.Info) {
	// encode body data
	var newDataTodo struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	err := json.NewDecoder(r.Body).Decode(&newDataTodo)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	todo := todos.Todo{
		ID:          mux.Vars(r)["id"],
		Title:       newDataTodo.Title,
		Description: newDataTodo.Description,
		IsDone:      false, // пока что так
	}

	if ok := todosDB.UserHaveThisTodo(userInfo, todo); !ok {
		http.Error(w, "Invalid todo id", http.StatusForbidden)
		return
	}

	_, err = todosDB.UpdateTodo(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handlerTodosDelete(w http.ResponseWriter, r *http.Request, userInfo user.Info) {
	var todo todos.Todo
	todo.ID = mux.Vars(r)["id"]

	if ok := todosDB.UserHaveThisTodo(userInfo, todo); !ok {
		http.Error(w, "Invalid todo id", http.StatusForbidden)
		return
	}

	_, err := todosDB.DeleteTodo(todo)
	if err != nil {
		http.Error(w, "Error delete todo", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// this func update todos status
func handlerTodosPatch(w http.ResponseWriter, r *http.Request, userInfo user.Info) {
	var todo todos.Todo
	var todoStatus struct {
		Status bool `json:"isDone"`
	}
	todo.ID = mux.Vars(r)["id"]

	err := json.NewDecoder(r.Body).Decode(&todoStatus)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	todo.IsDone = todoStatus.Status

	if ok := todosDB.UserHaveThisTodo(userInfo, todo); !ok {
		http.Error(w, "Invalid todo id", http.StatusForbidden)
		return
	}
	_, err = todosDB.UpdateTodoStatus(todo)
	if err != nil {
		http.Error(w, "Error update todo status", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}
