package repository

import (
	"database/sql"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
)

type FilesPostgres struct {
	db *sql.DB
}

func NewFilesPostgres(db *sql.DB) *FilesPostgres {
	return &FilesPostgres{
		db: db,
	}
}

func (r *FilesPostgres) Create(file domain.File) error {
	_, err := r.db.Exec("UPDATE products SET image=$1 WHERE id = $2", file.URL, file.ProductId)

	return err
}

func (r *FilesPostgres) GetProductImage(productId string) (string, error) {
	var productImage string
	rows, err := r.db.Query("SELECT image FROM products WHERE id = $1", productId)
	if err != nil {
		return productImage, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&productImage); err != nil {
			return productImage, err
		}
	}

	return productImage, rows.Err()
}
