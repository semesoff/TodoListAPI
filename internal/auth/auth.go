package auth

import (
	"errors"
	jwt2 "github.com/golang-jwt/jwt"
	"net/http"
	dbUser "todo-list-api/internal/db/user"
	"todo-list-api/internal/jwt"
	"todo-list-api/internal/models/user"
)

func HandleLoginWithToken(w http.ResponseWriter, r *http.Request) (user.Info, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return user.Info{}, errors.New("cookie is empty")
	}

	parsedToken, err := jwt.ValidateJWT(cookie.Value)
	if err != nil {
		return user.Info{}, err
	}

	claims, ok := parsedToken.Claims.(jwt2.MapClaims)
	if !ok || !parsedToken.Valid {
		return user.Info{}, errors.New("cookie is invalid")
	}

	email := claims["email"].(string)
	fullUserData, err := dbUser.GetUserByEmail(user.Login{
		Email:    email,
		Password: "",
	})

	if err != nil {
		return user.Info{}, err
	}

	return user.Info{ID: fullUserData.ID, Username: fullUserData.Username, Email: fullUserData.Email}, nil
}
