package user

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt string
}

type Login struct {
	Email    string
	Password string
}

type Register struct {
	Username string
	Email    string
	Password string
}

type Info struct {
	ID       int
	Username string
	Email    string
}
