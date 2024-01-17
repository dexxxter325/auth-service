package service

import (
	"CRUD_API"
	"CRUD_API/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Product interface {
	Create(name, description string) (CRUD_API.Products, error)
	ReadAll() ([]CRUD_API.Products, error)
	ReadById(id int) (CRUD_API.Products, error)
	Update(name, description string, id int) (CRUD_API.Products, error)
	Delete(id int) error
}

type Authorization interface {
	CreateUser(user CRUD_API.User) (int, error)                                       //аутентификация
	GenerateAccessToken(username, password string, hashPassword bool) (string, error) //авторизация(тк токен будет использован для доступа к данным)
	ParseAccessToken(token string) (int, error)
	GenerateRefreshToken(int) (string, error)
	ParseRefreshToken(refreshToken string) (int, error)
	GenerateNewTokenPair(refreshToken string) (string, string, error)
	/*
		Для начала система запрашивает логин, пользователь его указывает, система распознает его как существующий - это идентификация.
		После этого Google просит ввести пароль, пользователь его вводит, и система соглашается, что пользователь, похоже, действительно настоящий, раз пароль совпал, - это аутентификация.
		Скорее всего, Google дополнительно спросит еще и одноразовый код из SMS или приложения. Если пользователь и его правильно введет, то система окончательно согласится с тем, что он настоящий владелец аккаунта, - это двухфакторная аутентификация.
		После этого система предоставит пользователю право читать письма в его почтовом ящике и все в таком духе - это авторизация.*/
}
type Service struct {
	Product
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		NewServiceProduct(repos.Product),
		NewServiceAuthorization(repos.Authorization),
	}
}
