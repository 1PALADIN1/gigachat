package service

import (
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
	"github.com/gorilla/websocket"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	SignUpUser(user entity.User) (int, error)
	GenerateToken(username, password string) (string, int, error)
	ParseToken(token string) (int, error)
}

type User interface {
	GetUserById(id int) (entity.User, error)
	FindUserByName(filter string, currentUserId int) ([]entity.User, error)
}

type Chat interface {
	GetOrCreateChat(chat entity.Chat) (int, error)
	GetAllChats(userId int) ([]entity.ChatResponse, error)
}

type Message interface {
	AddMessageToChat(userId, chatId int, message string) (entity.Message, error)
	GetAllMessages(chatId int) ([]entity.Message, error)
}

type UserConnection interface {
	AddUserInActiveList(userId int, connection *websocket.Conn)
	RemoveUserFromActiveList(userId int)
	NotifyActiveUsers(notification NotificationMessage) error
	CloseAllConnections()
}

type Notification interface {
	NewMessageNotification(messageType int, message entity.ResponseMessage) (NotificationMessage, error)
}

type Service struct {
	Authorization
	User
	Chat
	Message
	UserConnection
	Notification
}

type ServiceConfig struct {
	Auth struct {
		Addr        string
		ConnTimeout int
	}
	App struct {
		MinSearchSymbols int
	}
}

func NewService(repo *repository.Repository, config ServiceConfig) *Service {
	return &Service{
		Authorization:  NewAuthService(config.Auth.Addr, time.Duration(config.Auth.ConnTimeout)*time.Second),
		User:           NewUserService(repo.User, config.App.MinSearchSymbols),
		Chat:           NewChatService(repo.Chat, repo.User),
		Message:        NewMessageService(repo.Message),
		UserConnection: NewUserConnectionService(repo.Chat),
		Notification:   NewNotificationService(repo.Chat),
	}
}
