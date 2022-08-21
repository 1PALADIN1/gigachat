package postgres

import (
	"fmt"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db}
}

func (r *UserPostgres) GetUserById(id int) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, username FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}
