package postgres

import (
	"fmt"
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	//tables
	usersTable      = "users"
	chatsTable      = "chats"
	usersChatsTable = "users_chats"
	messagesTable   = "messages"

	//db
	connectInterval = 2 * time.Second
)

func NewDB(dsn string, connectionTimeout float64) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	logger.LogInfo("Trying to connect to database...")
	startTime := time.Now()
	var succeed bool
	for time.Since(startTime).Seconds() < connectionTimeout {
		err = db.Ping()
		if err == nil {
			succeed = true
			break
		}

		logger.LogInfo(fmt.Sprintf("failed connect: %s", err.Error()))
		time.Sleep(connectInterval)
	}

	if !succeed {
		return nil, err
	}

	logger.LogInfo("Connected to database!")
	return db, err
}
