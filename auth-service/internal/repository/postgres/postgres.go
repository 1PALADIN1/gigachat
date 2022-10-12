package postgres

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	//tables
	usersTable = "users"

	//db
	connectInterval = 2 * time.Second
)

func NewDB(dsn string, connectionTimeout float64) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

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

	log.Printf("Connected to database!")
	return db, err
}
