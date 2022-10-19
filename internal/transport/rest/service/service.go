package service

import (
	"context"

	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/repository"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/storage"
	audit "github.com/AndrewMislyuk/go-shop-logger/pkg/domain"
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

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type Service struct {
	User
	ProductsList
	Files
	AuditClient
}

func NewService(repos *repository.Repository, auditClient AuditClient, storage storage.Provider) *Service {
	return &Service{
		User:         NewAuthService(repos.Authorization),
		ProductsList: NewProductsListService(repos.ProductsList, storage, auditClient),
		Files:        NewFileService(repos.Files, storage),
	}
}
