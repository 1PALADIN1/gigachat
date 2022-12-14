package repository

import (
	"github.com/1PALADIN1/gigachat_server/auth/internal/entity"
	"github.com/1PALADIN1/gigachat_server/auth/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username, password string) (entity.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db),
	}
}
