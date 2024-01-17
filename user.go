package CRUD_API

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"` //login
	Password string `json:"password"`
}
