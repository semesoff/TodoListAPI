package register

import (
	"html/template"
	"log"
	"net/http"
	dbUser "todo-list-api/internal/db/user"
	"todo-list-api/internal/jwt"
	"todo-list-api/internal/models/user"
)

func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderRegisterPage(w)
		return
	}
	HandlerRegisterPost(w, r)
}

func renderRegisterPage(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("web/templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func HandlerRegisterPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := jwt.HashPassword(r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userRegisterData := user.Register{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: hashedPassword,
	}

	// user already in database or no
	if _, ok := dbUser.GetUserByEmail(user.Login{
		Email:    userRegisterData.Email,
		Password: userRegisterData.Password},
	); ok == nil {
		http.Error(w, "User already exist.", http.StatusForbidden)
		return
	}

	if err := dbUser.AddUser(userRegisterData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fullUserData, err := dbUser.GetUserByEmail(
		user.Login{
			Email:    userRegisterData.Email,
			Password: userRegisterData.Password})

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

	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
