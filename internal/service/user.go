package service

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	mx          sync.Mutex
	activeUsers map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
)

func StartUserSession(connection *websocket.Conn) {
	addUserInList(connection)
	go handleUserMessages(connection)
}

func handleUserMessages(connection *websocket.Conn) {
	defer connection.Close()
	defer removeUserFromList(connection)

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			log.Println("Connection closed ", connection.RemoteAddr())
			break
		}

		sendMessageToAllUsers(messageType, string(message))
	}
}

func addUserInList(connection *websocket.Conn) {
	mx.Lock()

	activeUsers[connection] = true
	log.Println("Set user active", connection.RemoteAddr(), "active users:", len(activeUsers))

	mx.Unlock()
}

func removeUserFromList(connection *websocket.Conn) {
	mx.Lock()

	_, ok := activeUsers[connection]
	if ok {
		delete(activeUsers, connection)
		log.Println("Remove user from active list", connection.RemoteAddr(), "active users:", len(activeUsers))
	}

	mx.Unlock()
}

func sendMessageToAllUsers(messageType int, message string) {
	log.Println("Message type", messageType, "-> message:", message)
	mx.Lock()

	for conn := range activeUsers {
		conn.WriteMessage(messageType, []byte(message))
	}

	mx.Unlock()
}
