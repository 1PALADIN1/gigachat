package service

import (
	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository"
)

type NotificationMessage struct {
	MessageType int
	Users       []int
	Command     string
	Payload     interface{}
}

type NotificationService struct {
	chatRepo repository.Chat
}

const (
	messageCommand = "msg"
)

func NewNotificationService(chatRepo repository.Chat) *NotificationService {
	return &NotificationService{
		chatRepo: chatRepo,
	}
}

func (s *NotificationService) NewMessageNotification(messageType int, message entity.ResponseMessage) (NotificationMessage, error) {
	notificationMessage := NotificationMessage{}
	userIds, err := s.chatRepo.GetUserIdsByChatId(message.ChatId)
	if err != nil {
		return notificationMessage, err
	}

	notificationMessage.MessageType = messageType
	notificationMessage.Users = userIds
	notificationMessage.Command = messageCommand
	notificationMessage.Payload = message

	return notificationMessage, nil
}
