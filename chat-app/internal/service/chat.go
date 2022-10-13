package service

import (
	"fmt"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
)

type ChatService struct {
	chatRepo repository.Chat
	userRepo repository.User
}

func NewChatService(chatRepo repository.Chat, userRepo repository.User) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
		userRepo: userRepo,
	}
}

// GetOrCreateChat ищет или создаёт новый чат для пользователей
func (s *ChatService) GetOrCreateChat(chat entity.Chat) (int, error) {
	for _, userId := range chat.UserIds {
		if _, err := s.userRepo.GetUserById(userId); err != nil {
			return 0, fmt.Errorf("user with id=%d is not found", userId)
		}
	}

	chatId, found, err := s.chatRepo.GetChatIdByUsers(chat.UserIds)
	if err != nil {
		return 0, err
	}

	// чат нашёлся
	if found {
		return chatId, nil
	}

	chatId, err = s.chatRepo.CreateChat(chat)
	if err != nil {
		return 0, err
	}

	return chatId, nil
}

// GetAllChats получение списка чатов пользователя
func (s *ChatService) GetAllChats(userId int) ([]entity.ChatResponse, error) {
	return s.chatRepo.GetAllChats(userId)
}
