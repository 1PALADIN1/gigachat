package service

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
	"github.com/gorilla/websocket"
)

type UserService struct {
	mx               sync.Mutex
	activeUsers      map[*websocket.Conn]bool
	repo             repository.User
	minSearchSymbols int
}

func NewUserService(repo repository.User, minSearchSymbols int) *UserService {
	return &UserService{
		repo:             repo,
		activeUsers:      make(map[*websocket.Conn]bool),
		minSearchSymbols: minSearchSymbols,
	}
}

// Ищет пользователя по id
// Возвращает структуру пользователя в случае успеха
func (s *UserService) GetUserById(id int) (entity.User, error) {
	return s.repo.GetUserById(id)
}

// Ищет пользователей по username (исключаем текущего пользователя)
func (s *UserService) FindUserByName(filter string, currentUserId int) ([]entity.User, error) {
	if len(strings.TrimSpace(filter)) < s.minSearchSymbols {
		return nil, fmt.Errorf("requires min %d symbols to perfom searching", s.minSearchSymbols)
	}

	return s.repo.FindUserByName(filter, currentUserId)
}

func (s *UserService) AddUserInActiveList(connection *websocket.Conn) {
	defer s.mx.Unlock()
	s.mx.Lock()

	s.activeUsers[connection] = true
	log.Println("Set user active", connection.RemoteAddr(), "active users:", len(s.activeUsers))
}

func (s *UserService) RemoveUserFromActiveList(connection *websocket.Conn) {
	defer s.mx.Unlock()
	s.mx.Lock()

	_, ok := s.activeUsers[connection]
	if ok {
		delete(s.activeUsers, connection)
		log.Println("Remove user from active list", connection.RemoteAddr(), "active users:", len(s.activeUsers))
	}
}

func (s *UserService) SendMessageToAllUsers(messageType int, message []byte) {
	defer s.mx.Unlock()
	s.mx.Lock()

	for conn := range s.activeUsers {
		conn.WriteMessage(messageType, message)
	}
}
