package postgres

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host              string
	Port              int
	User              string
	Password          string
	DBName            string
	SSLMode           string
	ConnectionTimeout int
}

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

	log.Printf("Trying to connect to databse...")
	startTime := time.Now()
	var succeed bool
	for time.Since(startTime).Seconds() < connectionTimeout {
		err = db.Ping()
		if err == nil {
			succeed = true
			break
		}

		log.Printf("failed connect: %s", err.Error())
		time.Sleep(connectInterval)
	}

	if !succeed {
		return nil, err
	}

	log.Printf("Connected to databse!")
	return db, err
}
