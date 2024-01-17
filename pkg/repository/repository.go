package repository

import (
	"CRUD_API"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Product interface {
	Create(name, description string) (CRUD_API.Products, error)
	ReadAll() ([]CRUD_API.Products, error)
	ReadById(id int) (CRUD_API.Products, error)
	Update(name, description string, id int) (CRUD_API.Products, error)
	Delete(id int) error
}
type Authorization interface {
	CreateUser(user CRUD_API.User) (int, error)
	GetUser(username, password string) (CRUD_API.User, error)
	GetUserById(id int) (CRUD_API.User, error)
}
type Repository struct {
	Product
	Authorization
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Product:       NewProductPostgres(db),
		Authorization: NewAuthorizationPostgres(db),
	}
}
