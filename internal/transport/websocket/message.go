package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type requestMessage struct {
	Text   string `json:"message"`
	ChatId int    `json:"chat_id"`
}

type responseMessage struct {
	SendTime     string `json:"send_time"`
	Text         string `json:"text"`
	ChatId       int    `json:"chat_id"`
	responseUser `json:"user"`
}

type responseUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

const timeFormat = "2006-01-02 15:04"

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

		var reqMessage requestMessage
		if err := json.Unmarshal(message, &reqMessage); err != nil {
			log.Println("error parsing request:", err.Error())
			continue
		}

		resMessage, err := h.service.AddMessageToChat(user.Id, reqMessage.ChatId, reqMessage.Text)
		if err != nil {
			log.Println(err)
			continue
		}

		rawMsg := responseMessage{
			responseUser: responseUser{
				Id:       user.Id,
				Username: user.Username,
			},
			SendTime: resMessage.SendTime.Format(timeFormat),
			Text:     reqMessage.Text,
			ChatId:   1,
		}

		jsonMsg, err := json.Marshal(&rawMsg)
		if err != nil {
			log.Println("error converting message:", err)
			continue
		}

		log.Println("Message type", messageType, "-> message:", string(jsonMsg))
		h.service.SendMessageToAllUsers(messageType, jsonMsg)
	}
}
