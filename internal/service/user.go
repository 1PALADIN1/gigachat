package service

import (
	"log"
	"sync"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
	"github.com/gorilla/websocket"
)

type UserService struct {
	mx          sync.Mutex
	activeUsers map[*websocket.Conn]bool
	repo        repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo:        repo,
		activeUsers: make(map[*websocket.Conn]bool),
	}
}

func (s *UserService) GetUserById(id int) (entity.User, error) {
	return s.repo.GetUserById(id)
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
