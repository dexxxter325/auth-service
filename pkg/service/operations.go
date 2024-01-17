package service

import (
	"CRUD_API"
	"CRUD_API/pkg/repository"
)

type ServiceProduct struct {
	repo repository.Product
}

func NewServiceProduct(repo repository.Product) *ServiceProduct {
	return &ServiceProduct{repo: repo}
}

func (s *ServiceProduct) Create(name, description string) (CRUD_API.Products, error) { //реализует интерфейс Product.
	return s.repo.Create(name, description)
}
func (s *ServiceProduct) ReadAll() ([]CRUD_API.Products, error) {
	return s.repo.ReadAll()
}
func (s *ServiceProduct) ReadById(id int) (CRUD_API.Products, error) {
	return s.repo.ReadById(id)
}

func (s *ServiceProduct) Update(name, description string, id int) (CRUD_API.Products, error) {
	return s.repo.Update(name, description, id)
}

func (s *ServiceProduct) Delete(id int) error {
	return s.repo.Delete(id)
}
