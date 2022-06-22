package repository

import (
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"

	"database/sql"
)

type Authorization interface {
	CreateUser(user domain.UserSignUp) (string, error)
	GetUser(email, password string) (domain.User, error)
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

type Repository struct {
	Authorization
	CategoriesList
	ProductsList
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization:  NewAuthPostgres(db),
		CategoriesList: NewCategoriesListPostgres(db),
		ProductsList:   NewProductsListPostgres(db),
	}
}
