package mysql

import (
	"fmt"

	"github.com/1PALADIN1/gigachat_server/log/internal/entity"
	"github.com/jmoiron/sqlx"
)

type LogMySQL struct {
	db *sqlx.DB
}

func NewLogMySQL(db *sqlx.DB) *LogMySQL {
	return &LogMySQL{
		db: db,
	}
}

// InsertLog вставляет новую запись в таблицу с логами
func (r *LogMySQL) InsertLog(entry entity.Log) error {
	query := fmt.Sprintf("INSERT INTO %s (level, source, message) VALUES (?, ?, ?)", logsTable)
	_, err := r.db.Exec(query, entry.Level, entry.Source, entry.Message)
	return err
}
