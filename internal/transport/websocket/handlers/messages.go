package handlers

import (
	"log"

	"github.com/1PALADIN1/gigachat_server/internal/service"
	"github.com/gorilla/websocket"
)

func HandleUserMessages(connection *websocket.Conn) {
	defer connection.Close()
	defer service.RemoveUserFromActiveList(connection)

	log.Println("Client", connection.RemoteAddr(), "connected!")
	service.AddUserInActiveList(connection)

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			log.Println("Connection closed ", connection.RemoteAddr())
			break
		}

		service.SendMessageToAllUsers(messageType, string(message))
	}
}
