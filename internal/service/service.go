package service

import (
	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
)

type Authorization interface {
	SignUpUser(user entity.User) (int, error)
}

type Service struct {
	Authorization
}

type AuthConfig struct {
	SigningKey       string
	PasswordHashSalt string
}

func NewService(repo *repository.Repository, authConfig AuthConfig) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, authConfig.SigningKey, authConfig.PasswordHashSalt),
	}
}
