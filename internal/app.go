// Стартовая точка приложения.
// Здесь создаются основные сущности, строятся зависимости.
package internal

import (
	"log"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/service"

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
		TokenTTL         int
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

	repo := repository.NewRepository(db)
	service := service.NewService(repo, service.AuthConfig{
		SigningKey:       config.Auth.SigningKey,
		PasswordHashSalt: config.Auth.PasswordHashSalt,
		TokenTTL:         config.Auth.TokenTTL,
	})

	server := setupServer(config, service)
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

func setupServer(config *Config, service *service.Service) *Server {
	mux := http.NewServeMux()
	wsHandler := websocket.NewHandler()
	wsHandler.SetupRoutes(mux)

	restHandler := rest.NewHandler(service)
	restHandler.SetupRoutes(mux)

	return NewServer(ServerConfig{
		Port:           config.Server.Port,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, //1MB
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
	})
}
