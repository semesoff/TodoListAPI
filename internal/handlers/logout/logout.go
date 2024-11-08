package logout

import "net/http"

func HandlerLogout(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("token"); err != nil {
		http.Error(w, "Error logout", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
