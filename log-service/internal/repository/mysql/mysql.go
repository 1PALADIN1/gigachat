package mysql

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	//tables
	logsTable = "logs"

	//db
	connectInterval = 2 * time.Second
)

// NewDB создаёт подключение к MySQL и тестирует соединение
func NewDB(dsn string, connectionTimeout float64) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	log.Println("Trying to connect to database...")
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

	log.Println("Connected to database!")
	return db, err
}
