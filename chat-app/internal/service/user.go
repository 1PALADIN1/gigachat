package service

import (
	"fmt"
	"strings"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
)

type UserService struct {
	repo             repository.User
	minSearchSymbols int
}

func NewUserService(repo repository.User, minSearchSymbols int) *UserService {
	return &UserService{
		repo:             repo,
		minSearchSymbols: minSearchSymbols,
	}
}

// GetUserById ищет пользователя по id
// Возвращает структуру пользователя в случае успеха
func (s *UserService) GetUserById(id int) (entity.User, error) {
	return s.repo.GetUserById(id)
}

// FindUserByName ищет пользователей по username (исключаем текущего пользователя)
func (s *UserService) FindUserByName(filter string, currentUserId int) ([]entity.User, error) {
	if len(strings.TrimSpace(filter)) < s.minSearchSymbols {
		return nil, fmt.Errorf("requires min %d symbols to perfom searching", s.minSearchSymbols)
	}

	return s.repo.FindUserByName(filter, currentUserId)
}
