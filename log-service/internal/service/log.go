package service

import (
	"fmt"
	"log"

	"github.com/1PALADIN1/gigachat_server/log/internal/entity"
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
	log.Printf("%s[%s]:\t%s", logLevel, source, message)
	err := s.logRepo.InsertLog(entity.Log{
		Level:   logLevel,
		Source:  source,
		Message: message,
	})

	if err != nil {
		err = fmt.Errorf("error inserting log record: %s", err.Error())
		log.Println(err)
		return err
	}

	return nil
}
