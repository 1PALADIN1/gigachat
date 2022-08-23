package service

import (
	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
	"github.com/gorilla/websocket"
)

type Authorization interface {
	SignUpUser(user entity.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type User interface {
	GetUserById(id int) (entity.User, error)
	FindUserByName(filter string, currentUserId int) ([]entity.User, error)
	AddUserInActiveList(connection *websocket.Conn)
	RemoveUserFromActiveList(connection *websocket.Conn)
	SendMessageToAllUsers(messageType int, message []byte)
}

type Chat interface {
	GetOrCreateChat(chat entity.Chat) (int, error)
	GetAllChats(userId int) ([]entity.Chat, error)
}

type Message interface {
	AddMessageToChat(userId, chatId int, message string) (entity.Message, error)
	GetAllMessages(chatId int) ([]entity.Message, error)
}

type Service struct {
	Authorization
	User
	Chat
	Message
}

type ServiceConfig struct {
	Auth struct {
		SigningKey       string
		PasswordHashSalt string
		TokenTTL         int
	}
	App struct {
		MinSearchSymbols int
	}
}

func NewService(repo *repository.Repository, config ServiceConfig) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, config.Auth.SigningKey, config.Auth.PasswordHashSalt, config.Auth.TokenTTL),
		User:          NewUserService(repo.User, config.App.MinSearchSymbols),
		Chat:          NewChatService(repo.Chat, repo.User),
		Message:       NewMessageService(repo.Message),
	}
}
