package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/google/uuid"
)

type ProductsListPostgres struct {
	db *sql.DB
}

func NewProductsListPostgres(db *sql.DB) *ProductsListPostgres {
	return &ProductsListPostgres{
		db: db,
	}
}

func (r *ProductsListPostgres) Create(list domain.CreateProductInput) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var productId string
	row, err := tx.Prepare("INSERT INTO products(id, title, image, price, sale, sale_old_price, category_id, created_at) values($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		return "", err
	}

	defer row.Close()

	if err = row.QueryRow(uuid.NewString(), list.Title, list.Image, list.Price, list.Sale, list.SaleOldPrice, list.CategoryId, time.Now()).Scan(&productId); err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return productId, nil
}

func (r *ProductsListPostgres) GetAll() ([]domain.ProductsList, error) {
	rows, err := r.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]domain.ProductsList, 0)
	for rows.Next() {
		var product domain.ProductsList
		if err := rows.Scan(&product.Id, &product.Title, &product.Image, &product.Price, &product.Sale, &product.SaleOldPrice, &product.CategoryId, &product.CreatedAt); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *ProductsListPostgres) GetById(listId string) (domain.ProductsList, error) {
	rows, err := r.db.Query("SELECT * FROM products WHERE id = $1", listId)
	if err != nil {
		return domain.ProductsList{}, err
	}

	var product domain.ProductsList
	for rows.Next() {
		if err := rows.Scan(&product.Id, &product.Title, &product.Image, &product.Price, &product.Sale, &product.SaleOldPrice, &product.CategoryId, &product.CreatedAt); err != nil {
			return product, err
		}
	}

	return product, rows.Err()
}

func (r *ProductsListPostgres) Update(itemId string, input domain.UpdateProductInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Image != nil {
		setValues = append(setValues, fmt.Sprintf("image=$%d", argId))
		args = append(args, *input.Image)
		argId++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	if input.Sale != nil {
		setValues = append(setValues, fmt.Sprintf("sale=$%d", argId))
		args = append(args, *input.Sale)
		argId++
	}

	if input.SaleOldPrice != nil {
		setValues = append(setValues, fmt.Sprintf("sale_old_price=$%d", argId))
		args = append(args, *input.SaleOldPrice)
		argId++
	}

	if input.CategoryId != nil {
		setValues = append(setValues, fmt.Sprintf("category_id=$%d", argId))
		args = append(args, *input.CategoryId)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d", setQuery, argId)

	args = append(args, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ProductsListPostgres) Delete(itemId string) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", itemId)

	return err
}
