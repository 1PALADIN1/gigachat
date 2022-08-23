package service

import (
	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
)

type MessageService struct {
	messageRepo repository.Message
}

func NewMessageService(messageRepo repository.Message) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

// Добавляет сообщение в указанный чат
func (s *MessageService) AddMessageToChat(userId, chatId int, message string) (entity.Message, error) {
	return s.messageRepo.AddMessageToChat(userId, chatId, message)
}
