package entity

import (
	"errors"
	"fmt"
	"strings"
)

type Chat struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	UserIds     []int  `json:"user_ids,omitempty"`
}

var duplicateUserIdChecker = make(map[int]int)

func (c Chat) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return errors.New("chat title is not set")
	}

	if c.UserIds == nil || len(c.UserIds) < 2 {
		return errors.New("needs at least 2 users to create chat")
	}

	// проверяем дубликаты id в запросе
	for k := range duplicateUserIdChecker {
		delete(duplicateUserIdChecker, k)
	}

	for _, userId := range c.UserIds {
		duplicateUserIdChecker[userId]++

		if duplicateUserIdChecker[userId] > 1 {
			return fmt.Errorf("duplicate user id=%d in request", userId)
		}
	}

	return nil
}
