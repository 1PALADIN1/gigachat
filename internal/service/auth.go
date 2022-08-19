package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/1PALADIN1/gigachat_server/internal/entity"

	"github.com/1PALADIN1/gigachat_server/internal/repository"
)

type AuthService struct {
	authRepo         repository.Authorization
	signingKey       string
	passwordHashSalt string
}

func NewAuthService(authRepo repository.Authorization, signingKey, passwordHashSalt string) *AuthService {
	return &AuthService{
		authRepo:         authRepo,
		signingKey:       signingKey,
		passwordHashSalt: passwordHashSalt,
	}
}

// Регистрирует нового пользователя.
// Возвращает ID ползователя в БД в случае успеха. В противном случае - ошибку.
func (s *AuthService) SignUpUser(user entity.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.authRepo.CreateUser(user)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(s.passwordHashSalt)))
}
