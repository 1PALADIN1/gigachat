package entity

import (
	"errors"
	"strings"
)

type User struct {
	Id       int    `json:"-" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (u User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New("username is not set")
	}

	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password is not set")
	}

	return nil
}
