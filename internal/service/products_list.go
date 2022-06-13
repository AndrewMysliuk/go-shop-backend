package service

import (
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/repository"
)

type ProductsListService struct {
	repo repository.ProductsList
}

func NewProductsListService(repo repository.ProductsList) *ProductsListService {
	return &ProductsListService{
		repo: repo,
	}
}

func (s *ProductsListService) Create(list domain.CreateProductInput) (string, error) {
	return s.repo.Create(list)
}

func (s *ProductsListService) GetAll() ([]domain.ProductsList, error) {
	return s.repo.GetAll()
}

func (s *ProductsListService) GetById(listId string) (domain.ProductsList, error) {
	return s.repo.GetById(listId)
}

func (s *ProductsListService) Update(itemId string, input domain.UpdateProductInput) error {
	return s.repo.Update(itemId, input)
}

func (s *ProductsListService) Delete(itemId string) error {
	return s.repo.Delete(itemId)
}