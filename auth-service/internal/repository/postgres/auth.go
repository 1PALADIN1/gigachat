package postgres

import (
	"errors"
	"fmt"
	"log"

	"github.com/1PALADIN1/gigachat_server/auth/internal/entity"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

// CreateUser создаёт нового пользователя в БД
func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) VALUES ($1, $2) RETURNING id", usersTable)

	if r.db == nil {
		log.Println("DB is nil!")
		return 0, nil
	}

	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		log.Println("error creating new user:", err)
		return 0, errors.New("error creating new user")
	}

	return id, nil
}

// GetUser пытается получить пользователя по указанному логину и хешу пароля
func (r *AuthPostgres) GetUser(username, password string) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}
