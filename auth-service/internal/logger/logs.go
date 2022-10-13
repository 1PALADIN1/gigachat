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
	// log levels
	Info    = "INFO"
	Warning = "WARN"
	Error   = "ERROR"

	connectAttempts = 20
	connectInterval = 2 * time.Second
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

func init() {
	infoLog = log.New(os.Stdout, fmt.Sprintf("%s\t", Info), log.Ldate|log.Ltime)
	warningLog = log.New(os.Stdout, fmt.Sprintf("%s\t", Warning), log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, fmt.Sprintf("%s\t", Error), log.Ldate|log.Ltime|log.Lshortfile)
}

// Setup устанавливает настройки для сервиса логов и пробует подключиться к нему
func Setup(a, s string, timeout int) error {
	addr = a
	source = s
	connTimeout = time.Duration(timeout) * time.Second

	LogInfo("Trying to connect to LogService...")
	var err error
	connectCount := 0
	for {
		if err = testConnection(); err == nil {
			// успешно подключились
			LogInfo("Connected to LogService!")
			isConfigured = true
			return nil
		}

		LogInfo(fmt.Sprintf("failed to connect to LogService: %s", err.Error()))

		connectCount++
		if connectCount > connectAttempts {
			// сервис недоступен
			break
		}

		time.Sleep(connectInterval)
	}

	isConfigured = false
	return err
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

// testConnection тестирует соединение с сервисом логирования
func testConnection() error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("error connecting to logger service: %s", err.Error())
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	_, err = c.Ping(ctx, &logs.PingRequest{})
	if err != nil {
		return fmt.Errorf("error sending log message: %s", err.Error())
	}

	return nil
}
