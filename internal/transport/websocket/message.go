package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type userMessage struct {
	User     string `json:"user"`
	SendTime string `json:"send_time"`
	Text     string `json:"text"`
}

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
			log.Println("Connection closed ", connection.RemoteAddr())
			break
		}

		rawMsg := userMessage{
			User:     user.Username,
			SendTime: time.Now().UTC().Format("2006-01-02 15:04"),
			Text:     string(message),
		}

		jsonMsg, err := json.Marshal(&rawMsg)
		if err != nil {
			log.Println("Error converting message:", err)
			continue
		}

		log.Println("Message type", messageType, "-> message:", string(jsonMsg))
		h.service.SendMessageToAllUsers(messageType, jsonMsg)
	}
}
