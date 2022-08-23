package entity

import "time"

type Message struct {
	Id       int
	SendTime time.Time
	Text     string
	UserId   int
	ChatId   int
}
