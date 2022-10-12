package service

import (
	"log"

	"github.com/1PALADIN1/gigachat_server/log/internal/repository"
)

type LogService struct {
	logRepo repository.Log
}

// const (
// 	Info    = "INFO"
// 	Warning = "WARN"
// 	Error   = "ERROR"
// )

func NewLogService(logRepo repository.Log) *LogService {
	return &LogService{
		logRepo: logRepo,
	}
}

// // LogInfo логирование сообщений уровня Info
// func (s *LogService) LogInfo(message string) error {
// 	return s.log(Info, message)
// }

// // LogWarning логирование сообщений уровня Warning
// func (s *LogService) LogWarning(message string) error {
// 	return s.log(Warning, message)
// }

// // LogError логирование сообщений уровня Error
// func (s *LogService) LogError(message string) error {
// 	return s.log(Error, message)
// }

// Log логирование сообщений
func (s *LogService) Log(logLevel, message string) error {
	log.Printf("%s: %s", logLevel, message)
	//TODO: insert into DB
	return nil
}
