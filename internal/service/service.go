package service

import (
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/repository"
)

type User interface {
	CreateUser(user domain.UserSignUp) (string, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (string, error)
}

type CategoriesList interface {
	Create(list domain.CreateCategoryInput) (string, error)
	GetAll() ([]domain.CategoriesList, error)
	GetById(listId string) (domain.CategoriesList, error)
	Update(itemId string, input domain.UpdateCategoryInput) error
	Delete(itemId string) error
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
	CategoriesList
	ProductsList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:  NewAuthService(repos.Authorization),
		CategoriesList: NewCategoriesListService(repos.CategoriesList),
		ProductsList:   NewProductsListService(repos.ProductsList),
	}
}
