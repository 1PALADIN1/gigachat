package repository

import (
	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username, password string) (entity.User, error)
}

type User interface {
	GetUserById(id int) (entity.User, error)
}

type Chat interface {
	GetChatIdByUsers(userIds []int) (int, bool, error)
	CreateChat(chat entity.Chat) (int, error)
}

type Repository struct {
	Authorization
	User
	Chat
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db),
		User:          postgres.NewUserPostgres(db),
		Chat:          postgres.NewChatPostgres(db),
	}
}
