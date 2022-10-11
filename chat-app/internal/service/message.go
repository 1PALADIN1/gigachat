package service

import (
	"time"

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
	sendTime := time.Now().UTC()
	return s.messageRepo.AddMessageToChat(userId, chatId, message, sendTime)
}

// Получение всех сообщений из указанного чата
func (s *MessageService) GetAllMessages(chatId int) ([]entity.Message, error) {
	return s.messageRepo.GetAllMessages(chatId)
}
