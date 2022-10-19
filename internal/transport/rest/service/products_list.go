package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/transport/rest/repository"
	"github.com/AndrewMislyuk/go-shop-backend/pkg/storage"
	audit "github.com/AndrewMislyuk/go-shop-logger/pkg/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductsListService struct {
	repo        repository.ProductsList
	storage     storage.Provider
	auditClient AuditClient
}

func NewProductsListService(repo repository.ProductsList, storage storage.Provider, auditClient AuditClient) *ProductsListService {
	return &ProductsListService{
		repo:        repo,
		storage:     storage,
		auditClient: auditClient,
	}
}

func (s *ProductsListService) Create(list domain.CreateProductInput) (string, error) {
	productId := uuid.New().String()
	timestamp := time.Now()

	return s.repo.Create(list, productId, timestamp)
}

func (s *ProductsListService) GetAll() ([]domain.ProductsList, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return []domain.ProductsList{}, err
	}

	// Start Test gRPC Logger -------------------------------------------------
	if err := s.auditClient.SendLogRequest(context.TODO(), audit.LogItem{
		Action:    audit.ACTION_GET,
		Entity:    audit.ENTITY_PRODUCT,
		EntityID:  uuid.New().String(),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Product.GetAll",
		}).Error("failed to send log request:", err)
		return []domain.ProductsList{}, err
	}
	// End Test gRPC Logger   -------------------------------------------------

	return products, nil
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
