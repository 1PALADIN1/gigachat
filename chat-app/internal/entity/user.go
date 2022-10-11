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

type ResponseUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
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

func (u User) ToResponse() ResponseUser {
	return ResponseUser{
		Id:       u.Id,
		Username: u.Username,
	}
}
