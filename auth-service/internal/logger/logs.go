package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1PALADIN1/gigachat_server/auth/internal/logger/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	Info    = "INFO"
	Warning = "WARN"
	Error   = "ERROR"
)

var (
	isConfigured bool
	addr         string
	source       string
	connTimeout  time.Duration

	// log channels
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
)

func Setup(a, s string, timeout int) {
	addr = a
	source = s
	connTimeout = time.Duration(timeout) * time.Second

	infoLog = log.New(os.Stdout, fmt.Sprintf("%s\t", Info), log.Ldate|log.Ltime)
	warningLog = log.New(os.Stdout, fmt.Sprintf("%s\t", Warning), log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, fmt.Sprintf("%s\t", Error), log.Ldate|log.Ltime|log.Lshortfile)

	isConfigured = true
}

// LogInfo логирование сообщений уровня Info
func LogInfo(message string) {
	infoLog.Println(message)
	pushLog(Info, message)
}

// LogWarning логирование сообщений уровня Warning
func LogWarning(message string) {
	warningLog.Println(message)
	pushLog(Warning, message)
}

// LogError логирование сообщений уровня Error
func LogError(message string) {
	errorLog.Println(message)
	pushLog(Error, message)
}

// log логирование сообщений
func pushLog(logLevel, message string) {
	if !isConfigured {
		return
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errorLog.Printf("error connecting to logger service: %s\n", err.Error())
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	_, err = c.Log(ctx, &logs.LogRequest{
		LogLevel: logLevel,
		Message:  message,
		Source:   source,
	})

	if err != nil {
		errorLog.Printf("error sending log message: %s\n", err.Error())
	}
}
