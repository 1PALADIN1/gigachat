package entity

import "time"

const MessageTimeFormat = "2006-01-02 15:04"

type Message struct {
	Id       int       `db:"id"`
	SendTime time.Time `db:"send_date_time"`
	Text     string    `db:"message"`
	UserId   int       `db:"user_id"`
	ChatId   int       `db:"chat_id"`
}

type RequestMessage struct {
	Text   string `json:"message"`
	ChatId int    `json:"chat_id"`
}

type ResponseMessage struct {
	SendTime     string `json:"send_time"`
	Text         string `json:"text"`
	ChatId       int    `json:"chat_id"`
	ResponseUser `json:"user"`
}

func (m Message) ToResponse(user User) ResponseMessage {
	return ResponseMessage{
		ResponseUser: user.ToResponse(),
		SendTime:     m.SendTime.Format(MessageTimeFormat),
		Text:         m.Text,
		ChatId:       m.ChatId,
	}
}
