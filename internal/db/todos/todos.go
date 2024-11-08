package todos

import (
	"errors"
	"todo-list-api/internal/db/db"
	"todo-list-api/internal/models/todos"
	"todo-list-api/internal/models/user"
)

func AddTodo(todo todos.Todo) (todos.Todo, error) {
	DB := db.GetDB()
	err := DB.QueryRow(""+
		"INSERT INTO todos (user_id, title, description, is_done) VALUES ($1, $2, $3, $4) RETURNING id",
		todo.UserId, todo.Title, todo.Description, todo.IsDone).Scan(&todo.ID)
	if err != nil {
		return todos.Todo{}, err
	}
	return todo, nil
}

func GetAllTodos(userInfo user.Info) ([]todos.Todo, error) {
	DB := db.GetDB()
	rows, err := DB.Query("SELECT id, user_id, title, description, is_done FROM todos WHERE user_id = $1", userInfo.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listTodos := make([]todos.Todo, 0)

	for rows.Next() {
		var todo todos.Todo
		if err = rows.Scan(&todo.ID, &todo.UserId, &todo.Title, &todo.Description, &todo.IsDone); err != nil {
			return nil, err
		}
		listTodos = append(listTodos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listTodos, nil
}

func GetTodosWithPagination(userInfo user.Info, limit int, offset int) ([]todos.Todo, error) {
	DB := db.GetDB()
	rows, err := DB.Query(
		"SELECT id, title, description, is_done FROM todos WHERE user_id = $1 LIMIT $2 OFFSET $3",
		userInfo.ID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listTodos []todos.Todo
	for rows.Next() {
		var todo todos.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsDone); err != nil {
			return nil, err
		}
		listTodos = append(listTodos, todo)
	}

	return listTodos, nil
}

func GetCountTodos(userInfo user.Info) (int, error) {
	DB := db.GetDB()

	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM todos WHERE user_id = $1", userInfo.ID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func UpdateTodo(todo todos.Todo) (todos.Todo, error) {
	DB := db.GetDB()
	_, err := DB.Exec(
		"UPDATE todos SET title = $1, description = $2 WHERE id = $3",
		todo.Title, todo.Description, todo.ID)
	if err != nil {
		return todos.Todo{}, errors.New("failed update todo")
	}
	return todo, nil
}

func DeleteTodo(todo todos.Todo) (todos.Todo, error) {
	DB := db.GetDB()
	err := DB.QueryRow(
		"DELETE FROM todos WHERE id = $1 RETURNING user_id, title, description, is_done",
		todo.ID).Scan(&todo.UserId, &todo.Title, &todo.Description, &todo.IsDone)
	if err != nil {
		return todos.Todo{}, errors.New("failed delete todo")
	}
	return todo, nil
}

func UpdateTodoStatus(todo todos.Todo) (todos.Todo, error) {
	DB := db.GetDB()
	_, err := DB.Exec("UPDATE todos SET is_done = $1 WHERE id = $2", todo.IsDone, todo.ID)
	if err != nil {
		return todos.Todo{}, err
	}
	return todo, nil
}

func UserHaveThisTodo(userInfo user.Info, todo todos.Todo) bool {
	DB := db.GetDB()
	row := DB.QueryRow("SELECT id FROM todos WHERE id = $1 AND user_id = $2", todo.ID, userInfo.ID)

	var todoID int
	if row.Scan(&todoID) != nil {
		return false
	}
	return true
}
