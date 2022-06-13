package service

import (
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/repository"
)

type CategoriesListService struct {
	repo repository.CategoriesList
}

func NewCategoriesListService(repo repository.CategoriesList) *CategoriesListService {
	return &CategoriesListService{
		repo: repo,
	}
}

func (s *CategoriesListService) Create(list domain.CreateCategoryInput) (string, error) {
	return s.repo.Create(list)
}

func (s *CategoriesListService) GetAll() ([]domain.CategoriesList, error) {
	return s.repo.GetAll()
}

func (s *CategoriesListService) GetById(listId string) (domain.CategoriesList, error) {
	return s.repo.GetById(listId)
}

func (s *CategoriesListService) Update(itemId string, input domain.UpdateCategoryInput) error {
	return s.repo.Update(itemId, input)
}

func (s *CategoriesListService) Delete(itemId string) error {
	return s.repo.Delete(itemId)
}
