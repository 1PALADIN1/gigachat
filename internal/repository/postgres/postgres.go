package postgres

import (
	"fmt"
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
	connectInterval = 2
)

func NewDB(config Config) (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	log.Printf("Trying to connect to databse...")
	startTime := time.Now()
	var succeed bool
	for time.Since(startTime).Seconds() < float64(config.ConnectionTimeout) {
		err = db.Ping()
		if err == nil {
			succeed = true
			break
		}

		log.Printf("failed connect: %s", err.Error())
		time.Sleep(connectInterval * time.Second)
	}

	if !succeed {
		return nil, err
	}

	log.Printf("Connected to databse!")
	return db, err
}
