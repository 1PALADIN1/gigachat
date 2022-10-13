package service

import (
	"encoding/json"
	"fmt"

	"github.com/1PALADIN1/gigachat_server/internal/logger"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
	"github.com/gorilla/websocket"
)

type addUserInfo struct {
	userId     int
	connection *websocket.Conn
}

type notificationInfo struct {
	notification NotificationMessage
	message      []byte
}

type UserConnectionService struct {
	activeUsers map[int]*websocket.Conn
	chatRepo    repository.Chat

	// channels
	notifyUsers      chan notificationInfo
	addUser          chan addUserInfo
	removeUser       chan int
	closeConnections chan bool
}

func NewUserConnectionService(chatRepo repository.Chat) *UserConnectionService {
	srv := &UserConnectionService{
		activeUsers:      make(map[int]*websocket.Conn),
		chatRepo:         chatRepo,
		notifyUsers:      make(chan notificationInfo),
		addUser:          make(chan addUserInfo),
		removeUser:       make(chan int),
		closeConnections: make(chan bool),
	}

	go srv.listenChannels()
	return srv
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

	s.notifyUsers <- notificationInfo{
		notification: notification,
		message:      message,
	}

	return nil
}

// AddUserInActiveList добавляет пользователя в список онлайн
func (s *UserConnectionService) AddUserInActiveList(userId int, connection *websocket.Conn) {
	s.addUser <- addUserInfo{
		userId:     userId,
		connection: connection,
	}
}

// RemoveUserFromActiveList удаляет пользователя из списка онлайн
func (s *UserConnectionService) RemoveUserFromActiveList(userId int) {
	s.removeUser <- userId
}

// CloseAllConnections закрывает все пользовательские соединения
func (s *UserConnectionService) CloseAllConnections() {
	s.closeConnections <- true
}

func (s *UserConnectionService) listenChannels() {
	for {
		select {
		case info := <-s.notifyUsers:
			for _, userId := range info.notification.Users {
				if conn, ok := s.activeUsers[userId]; ok {
					conn.WriteMessage(info.notification.MessageType, info.message)
				}
			}

		case user := <-s.addUser:
			s.activeUsers[user.userId] = user.connection

		case userId := <-s.removeUser:
			_, ok := s.activeUsers[userId]
			if ok {
				delete(s.activeUsers, userId)
			}

		case <-s.closeConnections:
			for id, conn := range s.activeUsers {
				if err := conn.Close(); err != nil {
					logger.LogError(fmt.Sprintf("error closing user [%d] connection: %s\n", id, err.Error()))
				}
			}

			return
		}
	}
}
