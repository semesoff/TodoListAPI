package login

import (
	"html/template"
	"log"
	"net/http"
	"todo-list-api/internal/auth"
	dbUser "todo-list-api/internal/db/user"
	"todo-list-api/internal/jwt"
	"todo-list-api/internal/models/user"
)

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.HandleLoginWithToken(w, r); ok == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		renderLoginPage(w)
		return
	}
	handleLoginPost(w, r)
}

func renderLoginPage(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func handleLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userLoginData := user.Login{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	_, err := dbUser.GetUserByLogin(userLoginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	fullUserData, err := dbUser.GetUserByLogin(
		user.Login{
			Email:    userLoginData.Email,
			Password: userLoginData.Password})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// generate token and after send it to the client
	token, err := jwt.GenerateJWT(fullUserData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
