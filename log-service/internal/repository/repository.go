package repository

import (
	"github.com/1PALADIN1/gigachat_server/log/internal/repository/mysql"
	"github.com/jmoiron/sqlx"
)

type Log interface {
}

type Repository struct {
	Log
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Log: mysql.NewLogMySQL(db),
	}
}
