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
	var username string

	query := fmt.Sprintf(`INSERT INTO %s (message, send_date_time, chat_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id`, messagesTable)
	sendTime := time.Now().UTC()

	row := r.db.QueryRow(query, message, sendTime, chatId, userId)
	err := row.Scan(&messageId)
	if err != nil {
		return resMessage, err
	}

	getUserQuery := fmt.Sprintf(`SELECT username FROM %s WHERE id=$1`, usersTable)
	row = r.db.QueryRow(getUserQuery, userId)
	if err := row.Scan(&username); err != nil {
		return resMessage, err
	}

	resMessage.Id = messageId
	resMessage.Text = message
	resMessage.SendTime = sendTime
	resMessage.ChatId = chatId
	resMessage.UserId = userId
	resMessage.Username = username

	return resMessage, nil
}

// Получение всех сообщений из указанного чата
func (r *MessagePostgress) GetAllMessages(chatId int) ([]entity.Message, error) {
	var messages []entity.Message
	query := fmt.Sprintf(`SELECT m.id, m.send_date_time, m.message, m.user_id, m.chat_id, u.username FROM %s m
						  INNER JOIN %s u ON u.id=m.user_id WHERE chat_id=$1`,
		messagesTable, usersTable)
	if err := r.db.Select(&messages, query, chatId); err != nil {
		return nil, err
	}

	return messages, nil
}
