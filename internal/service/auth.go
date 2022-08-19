package service

import (
	"crypto/sha1"
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
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.authRepo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	addExpireAt := time.Duration(s.tokenTTL) * time.Second
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(addExpireAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Id,
	})

	return token.SignedString([]byte(s.signingKey))
}

// Генерирует хеш пароля пользователя с помощью алгоритма sha1
func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(s.passwordHashSalt)))
}
