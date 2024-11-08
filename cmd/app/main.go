package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"todo-list-api/config"
	"todo-list-api/internal/db/db"
	"todo-list-api/internal/handlers/home"
	"todo-list-api/internal/handlers/login"
	"todo-list-api/internal/handlers/logout"
	"todo-list-api/internal/handlers/register"
	"todo-list-api/internal/handlers/todos"
)

func main() {
	config.InitConfig()
	db.InitDB()
	cfg := config.GetConfig().Server
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.PathPrefix("/scripts/").Handler(http.StripPrefix("/scripts/", http.FileServer(http.Dir("web/scripts"))))
	r.HandleFunc("/", home.HandlerHome).Methods("GET")
	r.HandleFunc("/login", login.HandlerLogin)
	r.HandleFunc("/register", register.HandlerRegister)
	r.HandleFunc("/logout", logout.HandlerLogout)
	r.HandleFunc("/todos", todos.HandlerTodos).Methods("GET", "POST")
	r.HandleFunc("/todos/{id}", todos.HandlerTodos).Methods("DELETE", "PUT", "PATCH")
	serverAddress := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("Server is running on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
