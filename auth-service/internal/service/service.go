package service

import (
	"github.com/1PALADIN1/gigachat_server/auth/internal/entity"
	"github.com/1PALADIN1/gigachat_server/auth/internal/repository"
)

type Authorization interface {
	SignUpUser(user entity.User) (int, error)
	GenerateToken(username, password string) (string, int, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
}

type ServiceConfig struct {
	Auth struct {
		SigningKey       string
		PasswordHashSalt string
		TokenTTL         int
	}
}

func NewService(repo *repository.Repository, config ServiceConfig) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, config.Auth.SigningKey, config.Auth.PasswordHashSalt, config.Auth.TokenTTL),
	}
}
