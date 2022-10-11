// Стартовая точка приложения.
// Здесь создаются основные сущности, строятся зависимости.
package internal

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1PALADIN1/gigachat_server/internal/service"
	"github.com/gorilla/mux"

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
		DSN               string
		ConnectionTimeout int
	}
	App struct {
		MinSearchSymbols int
	}
}

func Run(config *Config) {
	db, err := setupDB(config)
	if err != nil {
		log.Fatalf("error connecting to database: %s", err.Error())
	}

	repo := repository.NewRepository(db)

	srvConfig := service.ServiceConfig{}
	srvConfig.Auth.SigningKey = config.Auth.SigningKey
	srvConfig.Auth.PasswordHashSalt = config.Auth.PasswordHashSalt
	srvConfig.Auth.TokenTTL = config.Auth.TokenTTL
	srvConfig.App.MinSearchSymbols = config.App.MinSearchSymbols
	service := service.NewService(repo, srvConfig)

	server := setupServer(config, service)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("ChatApp started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	stopServer(server, db, service)
}

func setupDB(config *Config) (*sqlx.DB, error) {
	db, err := postgres.NewDB(config.DB.DSN, float64(config.DB.ConnectionTimeout))
	if err != nil {
		return nil, err
	}

	return db, err
}

func setupServer(config *Config, service *service.Service) *Server {
	r := mux.NewRouter()
	wsHandler := websocket.NewHandler(service)
	wsHandler.SetupRoutes(r)

	restHandler := rest.NewHandler(service)
	restHandler.SetupRoutes(r)

	return NewServer(ServerConfig{
		Port:           config.Server.Port,
		Handler:        r,
		MaxHeaderBytes: 1 << 20, //1MB
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
	})
}

func stopServer(server *Server, db *sqlx.DB, service *service.Service) {
	log.Println("ChatApp shutting down")

	service.UserConnection.CloseAllConnections()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("error shutting down server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("error closing db: %s", err.Error())
	}
}
