package service

import "github.com/1PALADIN1/gigachat_server/log/internal/repository"

type Log interface {
	Log(logLevel, message, source string) error
}

type Service struct {
	Log
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Log: NewLogService(repo.Log),
	}
}
