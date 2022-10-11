package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/golang-jwt/jwt/v4"

	"github.com/1PALADIN1/gigachat_server/internal/repository"
)

type AuthService struct {
	authRepo         repository.Authorization
	signingKey       string
	passwordHashSalt string
	tokenTTL         int
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewAuthService(authRepo repository.Authorization, signingKey, passwordHashSalt string, tokenTTL int) *AuthService {
	return &AuthService{
		authRepo:         authRepo,
		signingKey:       signingKey,
		passwordHashSalt: passwordHashSalt,
		tokenTTL:         tokenTTL,
	}
}

// Регистрирует нового пользователя.
// Возвращает ID ползователя в БД в случае успеха. В противном случае - ошибку.
func (s *AuthService) SignUpUser(user entity.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.authRepo.CreateUser(user)
}

// Генерирует JWT-токен после успешной авторизации пользователя
func (s *AuthService) GenerateToken(username, password string) (string, int, error) {
	user, err := s.authRepo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", 0, err
	}

	addExpireAt := time.Duration(s.tokenTTL) * time.Second
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(addExpireAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Id,
	})

	res, err := token.SignedString([]byte(s.signingKey))

	return res, user.Id, err
}

// Валидирует токен и возвращает id пользователя в случае успеха
func (s *AuthService) ParseToken(token string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims is not of type *tokenClaims")
	}

	return claims.UserId, nil
}

// Генерирует хеш пароля пользователя с помощью алгоритма sha1
func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(s.passwordHashSalt)))
}
