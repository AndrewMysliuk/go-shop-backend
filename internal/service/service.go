package service

import (
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/repository"
)

type User interface {
	CreateUser(user domain.UserSignUp) (string, error)
	GenerateToken(email, password string) (string, error)
	GetMe(token string) (domain.User, error)
}

type ProductsList interface {
	Create(list domain.CreateProductInput) (string, error)
	GetAll() ([]domain.ProductsList, error)
	GetById(listId string) (domain.ProductsList, error)
	Update(itemId string, input domain.UpdateProductInput) error
	Delete(itemId string) error
}

type Service struct {
	User
	ProductsList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:         NewAuthService(repos.Authorization),
		ProductsList: NewProductsListService(repos.ProductsList),
	}
}
