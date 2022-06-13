package repository

import (
	"database/sql"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/google/uuid"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user domain.UserSignUp) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var userId string
	row, err := tx.Prepare("INSERT INTO users(id, name, surname, email, phone, password_hash, created_at) values($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	if err != nil {
		return "", err
	}

	defer row.Close()

	if err = row.QueryRow(uuid.NewString(), user.Name, user.Surname, user.Email, user.Phone, user.Password, time.Now()).Scan(&userId); err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (r *AuthPostgres) GetUser(email, password string) (string, error) {
	var userId string
	rows, err := r.db.Query("SELECT id FROM users WHERE email = $1 AND password_hash = $2", email, password)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		if err := rows.Scan(&userId); err != nil {
			return "", err
		}
	}

	return userId, rows.Err()
}
