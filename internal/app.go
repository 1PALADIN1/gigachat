// Стартовая точка приложения.
// Здесь создаются основные сущности, строятся зависимости.
package internal

import (
	"log"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/repository"

	"github.com/jmoiron/sqlx"

	"github.com/1PALADIN1/gigachat_server/internal/repository/postgres"

	"github.com/1PALADIN1/gigachat_server/internal/transport/rest"
	"github.com/1PALADIN1/gigachat_server/internal/transport/websocket"
)

type Config struct {
	Server struct {
		Port         int
		ReadTimeout  int
		WriteTimeout int
	}
	Auth struct {
		SigningKey       string
		PasswordHashSalt string
	}
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
}

func Run(config *Config) {
	db, err := setupDB(config)
	if err != nil {
		log.Fatalf("error connecting to database: %s", err.Error())
	}

	_ = repository.NewRepository(db)

	server := setupServer(config)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func setupDB(config *Config) (*sqlx.DB, error) {
	db, err := postgres.NewDB(postgres.Config{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		User:     config.DB.User,
		Password: config.DB.Password,
		DBName:   config.DB.DBName,
		SSLMode:  config.DB.SSLMode,
	})

	if err != nil {
		return nil, err
	}

	return db, err
}

func setupServer(config *Config) *Server {
	mux := http.NewServeMux()
	wsHandler := websocket.NewHandler()
	wsHandler.SetupRoutes(mux)

	restHandler := rest.NewHandler()
	restHandler.SetupRoutes(mux)

	return NewServer(ServerConfig{
		Port:           config.Server.Port,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, //1MB
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
	})
}
