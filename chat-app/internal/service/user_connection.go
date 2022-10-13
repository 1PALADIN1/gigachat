package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/1PALADIN1/gigachat_server/internal/logger"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
	"github.com/gorilla/websocket"
)

type UserConnectionService struct {
	mx          sync.Mutex
	activeUsers map[int]*websocket.Conn
	chatRepo    repository.Chat
}

func NewUserConnectionService(chatRepo repository.Chat) *UserConnectionService {
	return &UserConnectionService{
		activeUsers: make(map[int]*websocket.Conn),
		chatRepo:    chatRepo,
	}
}

type NotificationResponse struct {
	Command string      `json:"cmd"`
	Payload interface{} `json:"payload"`
}

// NotifyActiveUsers сообщает активным пользователям о каком-либо событии (пришло новое сообщение, добавили в чат и т.п.)
func (s *UserConnectionService) NotifyActiveUsers(notification NotificationMessage) error {
	notificationMessage := NotificationResponse{
		Command: notification.Command,
		Payload: notification.Payload,
	}

	message, err := json.Marshal(&notificationMessage)
	if err != nil {
		return err
	}

	defer s.mx.Unlock()
	s.mx.Lock()

	for _, userId := range notification.Users {
		if conn, ok := s.activeUsers[userId]; ok {
			conn.WriteMessage(notification.MessageType, message)
		}
	}

	return nil
}

// AddUserInActiveList добавляет пользователя в список онлайн
func (s *UserConnectionService) AddUserInActiveList(userId int, connection *websocket.Conn) {
	defer s.mx.Unlock()
	s.mx.Lock()
	s.activeUsers[userId] = connection
}

// RemoveUserFromActiveList удаляет пользователя в списка онлайн
func (s *UserConnectionService) RemoveUserFromActiveList(userId int) {
	defer s.mx.Unlock()
	s.mx.Lock()

	_, ok := s.activeUsers[userId]
	if ok {
		delete(s.activeUsers, userId)
	}
}

// CloseAllConnections закрывает все пользовательские соединения
func (s *UserConnectionService) CloseAllConnections() {
	defer s.mx.Unlock()
	s.mx.Lock()

	for id, conn := range s.activeUsers {
		if err := conn.Close(); err != nil {
			logger.LogError(fmt.Sprintf("error closing user [%d] connection: %s\n", id, err.Error()))
		}
	}
}
