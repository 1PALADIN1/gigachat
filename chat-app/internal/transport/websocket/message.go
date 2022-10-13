package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/logger"
	"github.com/gorilla/websocket"
)

// handleUserMessages обработчик входящих websocket соединений
func (h *Handler) handleUserMessages(connection *websocket.Conn, userId int) {
	user, err := h.service.GetUserById(userId)
	if err != nil {
		logger.LogError(fmt.Sprintf("user with id %d is not found", userId))
		return
	}
	defer h.service.RemoveUserFromActiveList(userId)

	h.service.AddUserInActiveList(userId, connection)
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}

		var reqMessage entity.RequestMessage
		if err := json.Unmarshal(message, &reqMessage); err != nil {
			logger.LogError(fmt.Sprintf("error parsing request: %s", err.Error()))
			continue
		}

		resultMessage, err := h.service.AddMessageToChat(user.Id, reqMessage.ChatId, reqMessage.Text)
		if err != nil {
			logger.LogError(fmt.Sprintf("error adding message to chat %d: %s", reqMessage.ChatId, err.Error()))
			continue
		}

		notification, err := h.service.NewMessageNotification(messageType, resultMessage.ToResponse())
		if err != nil {
			logger.LogError(fmt.Sprintf("error creating message notification: %s", err.Error()))
			continue
		}

		if err := h.service.UserConnection.NotifyActiveUsers(notification); err != nil {
			logger.LogError(fmt.Sprintf("error notify users: %s", err.Error()))
		}
	}
}
