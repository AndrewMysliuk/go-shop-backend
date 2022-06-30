package repository

import (
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"

	"database/sql"
)

type Authorization interface {
	CreateUser(user domain.UserSignUp) (string, error)
	GetUser(email, password string) (domain.User, error)
}

type ProductsList interface {
	Create(list domain.CreateProductInput, productId string, timestamp time.Time) (string, error)
	GetAll() ([]domain.ProductsList, error)
	GetById(listId string) (domain.ProductsList, error)
	Update(itemId string, input domain.UpdateProductInput) error
	Delete(itemId string) error
}

type Repository struct {
	Authorization
	ProductsList
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		ProductsList:  NewProductsListPostgres(db),
	}
}
