package service

import (
	"log"

	"github.com/1PALADIN1/gigachat_server/log/internal/repository"
)

type LogService struct {
	logRepo repository.Log
}

func NewLogService(logRepo repository.Log) *LogService {
	return &LogService{
		logRepo: logRepo,
	}
}

// Log логирование сообщений
func (s *LogService) Log(logLevel, message, source string) error {
	log.Printf("%s\t[%s]: %s", logLevel, source, message)
	//TODO: insert into DB
	return nil
}
