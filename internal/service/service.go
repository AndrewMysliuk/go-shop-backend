package service

import (
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/repository"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/storage"
)

//go:generate mockgen -source=service.go -destination=mock/mock.go

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

type Files interface {
	Upload(file domain.File) (string, error)
}

type Service struct {
	User
	ProductsList
	Files
}

func NewService(repos *repository.Repository, storage storage.Provider) *Service {
	return &Service{
		User:         NewAuthService(repos.Authorization),
		ProductsList: NewProductsListService(repos.ProductsList, storage),
		Files:        NewFileService(repos.Files, storage),
	}
}
