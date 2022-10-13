package mysql

import "github.com/jmoiron/sqlx"

type LogMySQL struct {
	db *sqlx.DB
}

func NewLogMySQL(db *sqlx.DB) *LogMySQL {
	return &LogMySQL{
		db: db,
	}
}

//TODO: db methods
