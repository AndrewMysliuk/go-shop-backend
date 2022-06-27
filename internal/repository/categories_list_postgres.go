package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/google/uuid"
)

type CategoriesListPostgres struct {
	db *sql.DB
}

func NewCategoriesListPostgres(db *sql.DB) *CategoriesListPostgres {
	return &CategoriesListPostgres{
		db: db,
	}
}

func (r *CategoriesListPostgres) Create(list domain.CreateCategoryInput) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var categoryId string
	row, err := tx.Prepare("INSERT INTO categories(id, name, created_at) values($1, $2, $3) RETURNING id")
	if err != nil {
		return "", err
	}

	defer row.Close()

	if err = row.QueryRow(uuid.NewString(), list.Name, time.Now()).Scan(&categoryId); err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return categoryId, nil
}

func (r *CategoriesListPostgres) GetAll() ([]domain.CategoriesList, error) {
	rows, err := r.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]domain.CategoriesList, 0)
	for rows.Next() {
		var category domain.CategoriesList
		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, rows.Err()
}

func (r *CategoriesListPostgres) GetById(listId string) (domain.CategoriesList, error) {
	rows, err := r.db.Query("SELECT * FROM categories WHERE id = $1", listId)
	if err != nil {
		return domain.CategoriesList{}, err
	}
	defer rows.Close()

	var category domain.CategoriesList
	for rows.Next() {
		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt); err != nil {
			return category, err
		}
	}

	return category, rows.Err()
}

func (r *CategoriesListPostgres) Update(itemId string, inp domain.UpdateCategoryInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *inp.Name)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE categories SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *CategoriesListPostgres) Delete(itemId string) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", itemId)

	return err
}
