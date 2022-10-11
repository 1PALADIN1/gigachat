package websocket

import (
	"encoding/json"
	"log"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/gorilla/websocket"
)

func (h *Handler) handleUserMessages(connection *websocket.Conn, userId int) {
	user, err := h.service.GetUserById(userId)
	if err != nil {
		log.Printf("user with id %d is not found\n", userId)
		return
	}
	defer h.service.RemoveUserFromActiveList(userId)

	log.Println("User", userId, "connected!")
	h.service.AddUserInActiveList(userId, connection)

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			log.Println("connection closed", connection.RemoteAddr())
			break
		}

		var reqMessage entity.RequestMessage
		if err := json.Unmarshal(message, &reqMessage); err != nil {
			log.Println("error parsing request:", err.Error())
			continue
		}

		resultMessage, err := h.service.AddMessageToChat(user.Id, reqMessage.ChatId, reqMessage.Text)
		if err != nil {
			log.Println(err)
			continue
		}

		notification, err := h.service.NewMessageNotification(messageType, resultMessage.ToResponse())
		if err != nil {
			log.Println(err)
			continue
		}

		if err := h.service.UserConnection.NotifyActiveUsers(notification); err != nil {
			log.Println(err)
		}
	}
}
