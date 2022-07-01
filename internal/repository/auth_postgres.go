package repository

import (
	"database/sql"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user domain.UserSignUp, dataId string, timestamp time.Time) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var userId string
	row, err := tx.Prepare("INSERT INTO users(id, name, surname, email, phone, role, password_hash, created_at) values($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			logrus.Fatalf("query failed: %v, unable to abort: %v", err, rb)
		}

		return "", err
	}

	defer row.Close()

	if err = row.QueryRow(dataId, user.Name, user.Surname, user.Email, user.Phone, user.Role, user.Password, timestamp).Scan(&userId); err != nil {
		if rb := tx.Rollback(); rb != nil {
			logrus.Fatalf("query failed: %v, unable to abort: %v", err, rb)
		}

		return "", err
	}

	return userId, tx.Commit()
}

func (r *AuthPostgres) GetUser(email, password string) (domain.User, error) {
	var userData domain.User
	rows, err := r.db.Query("SELECT * FROM users WHERE email = $1 AND password_hash = $2", email, password)
	if err != nil {
		return userData, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&userData.Id, &userData.Name, &userData.Surname, &userData.Email, &userData.Phone, &userData.Role, &userData.Password, &userData.CreatedAt); err != nil {
			return userData, err
		}
	}

	return userData, rows.Err()
}
