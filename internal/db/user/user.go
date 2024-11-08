package user

import (
	"database/sql"
	"errors"
	"todo-list-api/internal/db/db"
	"todo-list-api/internal/jwt"
	"todo-list-api/internal/models/user"
)

func AddUser(user user.Register) error {
	DB := db.GetDB()

	_, err := DB.Exec(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		user.Username, user.Email, user.Password)
	return err
}

func GetUserByLogin(userLoginData user.Login) (user.User, error) {
	DB := db.GetDB()

	rows, err := DB.Query("SELECT * FROM users WHERE email = $1", userLoginData.Email)
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
		}
	}(rows)

	if err != nil {
		return user.User{}, err
	}

	var fullUserData user.User
	for rows.Next() {
		err := rows.Scan(&fullUserData.ID, &fullUserData.Username, &fullUserData.Email, &fullUserData.Password, &fullUserData.CreatedAt)
		if jwt.CheckPassword(fullUserData.Password, userLoginData.Password) == nil {
			return fullUserData, nil
		}
		if err != nil {
		}
	}
	return user.User{}, errors.New("user is not found")
}

func GetUserByEmail(userLoginData user.Login) (user.User, error) {
	DB := db.GetDB()

	rows, err := DB.Query("SELECT * FROM users WHERE email = $1", userLoginData.Email)
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
		}
	}(rows)

	if err != nil {
		return user.User{}, err
	}

	var fullUserData user.User
	for rows.Next() {
		err := rows.Scan(&fullUserData.ID, &fullUserData.Username, &fullUserData.Email, &fullUserData.Password, &fullUserData.CreatedAt)
		if fullUserData.Email == userLoginData.Email {
			return fullUserData, nil
		}
		if err != nil {
		}
	}
	return user.User{}, errors.New("user is not found")
}
