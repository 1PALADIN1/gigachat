package service

import (
	"encoding/json"
	"log"
	"sync"

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

// Сообщает активным пользователям о каком-либо событии (пришло новое сообщение, добавили в чат и т.п.)
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

			log.Println("Notify", notification.MessageType, "-> message:", string(message)) //TODO
		}
	}

	return nil
}

// Добавляет пользователя в список онлайн
func (s *UserConnectionService) AddUserInActiveList(userId int, connection *websocket.Conn) {
	defer s.mx.Unlock()
	s.mx.Lock()

	s.activeUsers[userId] = connection
	log.Println("Set user active", connection.RemoteAddr(), "active users:", len(s.activeUsers))
}

// Удаляет пользователя в списка онлайн
func (s *UserConnectionService) RemoveUserFromActiveList(userId int) {
	defer s.mx.Unlock()
	s.mx.Lock()

	_, ok := s.activeUsers[userId]
	if ok {
		delete(s.activeUsers, userId)
		log.Println("Remove user from active list", userId, "active users:", len(s.activeUsers))
	}
}

// Закрывает все пользовательские соединения
func (s *UserConnectionService) CloseAllConnections() {
	defer s.mx.Unlock()
	s.mx.Lock()

	for id, conn := range s.activeUsers {
		if err := conn.Close(); err != nil {
			log.Printf("error closing user [%d] connection: %s\n", id, err.Error())
		}
	}
}
