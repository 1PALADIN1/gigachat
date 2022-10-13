package postgres

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ChatPostgres struct {
	db *sqlx.DB
}

func NewChatPostgres(db *sqlx.DB) *ChatPostgres {
	return &ChatPostgres{db}
}

// GetChatIdByUsers пытается найти существующий чат по указанному набору пользователей
// Возвращает id чата в случае успеха
func (r *ChatPostgres) GetChatIdByUsers(userIds []int) (int, bool, error) {
	if len(userIds) == 0 {
		return 0, false, errors.New("invalid amount of users")
	}

	var sb strings.Builder
	queryArgs := make([]interface{}, 0)
	for i := 0; i < len(userIds); i++ {
		sb.WriteString("$")
		sb.WriteString(strconv.Itoa(i + 1))
		queryArgs = append(queryArgs, userIds[i])

		// в последний элемент не записываем ", "
		if i < len(userIds)-1 {
			sb.WriteString(", ")
		}
	}

	query := fmt.Sprintf(`SELECT out_t.chat_id FROM %s out_t
						  INNER JOIN (SELECT chat_id, COUNT(*) AS total FROM %s
						  			  GROUP BY chat_id) in_t ON out_t.chat_id=in_t.chat_id
						  WHERE user_id in (%s)
						  GROUP BY out_t.chat_id, in_t.total
						  HAVING COUNT(out_t.chat_id)=in_t.total AND in_t.total=%d`,
		usersChatsTable, usersChatsTable, sb.String(), len(userIds))

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return 0, false, err
	}

	for rows.Next() {
		var chatId int
		err := rows.Scan(&chatId)

		// чат найден
		if err == nil {
			return chatId, true, nil
		}

		if err != nil {
			return 0, false, err
		}
	}

	return 0, false, rows.Err()
}

// CreateChat создаёт новый чат и возвращает id чата в случае успеха
func (r *ChatPostgres) CreateChat(chat entity.Chat) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, nil
	}

	var chatId int
	createChatQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", chatsTable)
	row := r.db.QueryRow(createChatQuery, chat.Title, chat.Description)
	if err := row.Scan(&chatId); err != nil {
		tx.Rollback()
		return 0, err
	}

	addUserInChatQuery := fmt.Sprintf("INSERT INTO %s (user_id, chat_id) VALUES ($1, $2)", usersChatsTable)
	for _, userId := range chat.UserIds {
		if _, err := r.db.Exec(addUserInChatQuery, userId, chatId); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return chatId, tx.Commit()
}

// GetAllChats получение списка чатов пользователя
func (r *ChatPostgres) GetAllChats(userId int) ([]entity.ChatResponse, error) {
	var chats []entity.ChatResponse
	query := fmt.Sprintf(`SELECT c.id, c.title, c.description, t2.send_date_time, t2.message, t2.user_id, u.username FROM %s us
						  INNER JOIN %s c ON c.id=us.chat_id
						  LEFT JOIN (
							SELECT m.id, m.user_id, m.message, m.send_date_time, t1.chat_id FROM %s m
							LEFT JOIN (
								SELECT m.chat_id, MAX(m.send_date_time) AS send_date_time FROM %s m
								GROUP BY m.chat_id
							) t1 ON m.chat_id=t1.chat_id
							WHERE t1.send_date_time=m.send_date_time
						  ) t2 ON us.chat_id=t2.chat_id
						  INNER JOIN users u ON u.id=t2.user_id
						  WHERE us.user_id=$1`,
		usersChatsTable, chatsTable, messagesTable, messagesTable)

	err := r.db.Select(&chats, query, userId)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

// GetUserIdsByChatId получение пользователей из указанного чата
func (r *ChatPostgres) GetUserIdsByChatId(chatId int) ([]int, error) {
	var userIds []int
	query := fmt.Sprintf(`SELECT user_id FROM %s
						  WHERE chat_id=$1`, usersChatsTable)

	err := r.db.Select(&userIds, query, chatId)
	if err != nil {
		return nil, err
	}

	return userIds, nil
}
