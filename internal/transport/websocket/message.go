package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/service"
	"github.com/gorilla/websocket"
)

type Message struct {
	User     string `json:"user"`
	SendTime string `json:"send_time"`
	Text     string `json:"text"`
}

func (h *Handler) handleUserMessages(connection *websocket.Conn) {
	defer service.RemoveUserFromActiveList(connection)

	log.Println("Client", connection.RemoteAddr(), "connected!")
	service.AddUserInActiveList(connection)

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			log.Println("Connection closed ", connection.RemoteAddr())
			break
		}

		rawMsg := Message{
			User:     connection.RemoteAddr().String(),
			SendTime: time.Now().Format("2006-02-01 15:04"),
			Text:     string(message),
		}

		jsonMsg, err := json.Marshal(&rawMsg)
		if err != nil {
			log.Println("Error converting message:", err)
			continue
		}

		log.Println("Message type", messageType, "-> message:", string(jsonMsg))
		service.SendMessageToAllUsers(messageType, jsonMsg)
	}
}
