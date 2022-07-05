package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/repository"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/storage"
	"github.com/google/uuid"
)

type ProductsListService struct {
	repo    repository.ProductsList
	storage storage.Provider
}

func NewProductsListService(repo repository.ProductsList, storage storage.Provider) *ProductsListService {
	return &ProductsListService{
		repo:    repo,
		storage: storage,
	}
}

func (s *ProductsListService) Create(list domain.CreateProductInput) (string, error) {
	productId := uuid.New().String()
	timestamp := time.Now()

	return s.repo.Create(list, productId, timestamp)
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
	product, err := s.GetById(itemId)
	if err != nil {
		return err
	}

	if product.Image != "" {
		filename := s.parseImageURL(product.Image)

		err := s.storage.Delete(context.Background(), filename)
		if err != nil {
			return err
		}
	}

	return s.repo.Delete(itemId)
}

func (s *ProductsListService) parseImageURL(url string) string {
	str := strings.Split(url, "/")

	return fmt.Sprintf("images/%s", str[len(str)-1])
}
