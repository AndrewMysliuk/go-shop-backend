package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/repository"
	"github.com/golang-jwt/jwt"
)

const (
	salt       = "superhashkey"
	signingKey = "L6e2h3e6gfE4ae93AZMfPLRg782Y"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	User domain.User `json:"data"`
}

type Auth struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *Auth {
	return &Auth{
		repo: repo,
	}
}

func (a *Auth) CreateUser(user domain.UserSignUp) (string, error) {
	user.Password = generatePasswordHash(user.Password)

	return a.repo.CreateUser(user)
}

func (a *Auth) GenerateToken(email, password string) (string, error) {
	user, err := a.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user,
	})

	return token.SignedString([]byte(signingKey))
}

func (a *Auth) GetMe(accessToken string) (domain.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return domain.User{}, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return domain.User{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return domain.User{}, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.User, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
