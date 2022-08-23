package postgres

import (
	"fmt"
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/jmoiron/sqlx"
)

type MessagePostgress struct {
	db *sqlx.DB
}

func NewMessagePostgress(db *sqlx.DB) *MessagePostgress {
	return &MessagePostgress{db}
}

// Добавляет сообщение в указанный чат
func (r *MessagePostgress) AddMessageToChat(userId, chatId int, message string) (entity.Message, error) {
	var resMessage entity.Message
	var messageId int

	query := fmt.Sprintf(`INSERT INTO %s (message, send_date_time, chat_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id`, messagesTable)
	sendTime := time.Now().UTC()

	row := r.db.QueryRow(query, message, sendTime, chatId, userId)
	err := row.Scan(&messageId)
	if err != nil {
		return resMessage, err
	}

	resMessage.Id = messageId
	resMessage.Text = message
	resMessage.SendTime = sendTime
	resMessage.ChatId = chatId
	resMessage.UserId = userId

	return resMessage, nil
}
