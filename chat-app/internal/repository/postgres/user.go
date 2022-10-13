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

// GetUserById ищет пользователя по указанному id
func (r *UserPostgres) GetUserById(id int) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf(`SELECT id, username FROM %s WHERE id=$1`, usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

// FindUserByName ищет пользователей по username (исключаем текущего пользователя)
func (r *UserPostgres) FindUserByName(filter string, currentUserId int) ([]entity.User, error) {
	var users []entity.User
	query := fmt.Sprintf(`SELECT id, username FROM %s 
						  WHERE username LIKE $1 AND id<>$2`, usersTable)
	if err := r.db.Select(&users, query, "%"+filter+"%", currentUserId); err != nil {
		return nil, err
	}

	return users, nil
}
