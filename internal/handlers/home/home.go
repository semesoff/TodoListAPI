package home

import (
	"html/template"
	"log"
	"net/http"
	"todo-list-api/internal/auth"
)

type PageVariables struct {
	IsLoggedIn bool
	UserName   string
	UserEmail  string
}

func HandlerHome(w http.ResponseWriter, r *http.Request) {
	pageVars := PageVariables{}
	if userInfo, ok := auth.HandleLoginWithToken(w, r); ok == nil {
		pageVars.IsLoggedIn = true
		pageVars.UserName = userInfo.Username
		pageVars.UserEmail = userInfo.Email
	}
	renderHomePage(w, pageVars)
}

func renderHomePage(w http.ResponseWriter, pageVars PageVariables) {
	tmpl, err := template.ParseFiles("web/templates/home.html")
	if err != nil {
		http.Error(w, "Error parsing html-template.", http.StatusInternalServerError)
		log.Fatalf("Error parsing html-template (home): %v", err)
	}
	if err := tmpl.Execute(w, pageVars); err != nil {
		log.Fatalf("Error executing html-template (home): %v", err)
	}
}
