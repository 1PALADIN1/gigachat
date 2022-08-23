package websocket

import (
	"encoding/json"
	"log"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/gorilla/websocket"
)

func (h *Handler) handleUserMessages(connection *websocket.Conn, userId int) {
	defer h.service.RemoveUserFromActiveList(connection)

	log.Println("User", userId, "connected!")
	h.service.AddUserInActiveList(connection)

	user, err := h.service.GetUserById(userId)
	if err != nil {
		log.Printf("user with id %d is not found\n", userId)
		return
	}

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

		respMessage := resultMessage.BuildMessageResponse(user)
		jsonMsg, err := json.Marshal(&respMessage)
		if err != nil {
			log.Println("error converting message:", err)
			continue
		}

		log.Println("Message type", messageType, "-> message:", string(jsonMsg))
		h.service.SendMessageToAllUsers(messageType, jsonMsg)
	}
}
