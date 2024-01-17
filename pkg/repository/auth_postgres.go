package repository

import (
	"CRUD_API"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type AuthorizationPostgres struct {
	db *pgxpool.Pool
}

func NewAuthorizationPostgres(db *pgxpool.Pool) *AuthorizationPostgres {
	return &AuthorizationPostgres{db: db}
}

func (r *AuthorizationPostgres) CreateUser(user CRUD_API.User) (int, error) {
	var err error
	request := "INSERT INTO users (name,username,password) VALUES ($1,$2,$3) RETURNING id"
	dorequest := r.db.QueryRow(context.Background(), request, user.Name, user.Username, user.Password)
	if err := dorequest.Scan(&user.ID); err != nil {
		return 0, fmt.Errorf("err in scan(func createUser):%s", err)
	}
	return user.ID, err
}

/*GetUser проверяем,есть ли по такому лог/пэссу пользователь в бд и берем его id для дальнейшего вшивания в токен.*/
func (r *AuthorizationPostgres) GetUser(username, password string) (CRUD_API.User, error) {
	var user CRUD_API.User
	var err error
	request := "SELECT id, username, password FROM users WHERE username=$1 AND password=$2"
	dorequest := r.db.QueryRow(context.Background(), request, username, password)
	if err := dorequest.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		fmt.Println("Error scanning:", err)
		return user, err
	}
	return user, err
}

func (r *AuthorizationPostgres) GetUserById(id int) (CRUD_API.User, error) {
	var user CRUD_API.User
	var err error
	request := "SELECT username,password FROM users WHERE id=$1 "
	dorequest := r.db.QueryRow(context.Background(), request, id)
	if err := dorequest.Scan(&user.Username, &user.Password); err != nil {
		log.Printf("err in doreq.scan:%s", err)
		return user, err
	}
	return user, err
}
