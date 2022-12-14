package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/1PALADIN1/gigachat_server/internal/logger"
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
		Addr        string
		ConnTimeout int
	}
	Log struct {
		Addr        string
		ConnTimeout int
		Source      string
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
	if err := logger.Setup(config.Log.Addr, config.Log.Source, config.Log.ConnTimeout); err != nil {
		logger.LogError(fmt.Sprintf("log service is unavailable: %s", err.Error()))
		os.Exit(1)
	}

	db, err := setupDB(config)
	if err != nil {
		logger.LogError(fmt.Sprintf("error connecting to database: %s", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.LogError(fmt.Sprintf("error closing db: %s", err.Error()))
		}
	}()

	repo := repository.NewRepository(db)

	srvConfig := service.ServiceConfig{}
	srvConfig.Auth.Addr = config.Auth.Addr
	srvConfig.Auth.ConnTimeout = config.Auth.ConnTimeout
	srvConfig.App.MinSearchSymbols = config.App.MinSearchSymbols
	service := service.NewService(repo, srvConfig)

	server := setupServer(config, service)
	go func() {
		if err := server.Start(); err != nil {
			logger.LogError(err.Error())
		}
	}()

	logger.LogInfo("ChatApp started")
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
	logger.LogInfo("ChatApp shutting down")
	service.UserConnection.CloseAllConnections()

	if err := server.Shutdown(context.Background()); err != nil {
		logger.LogError(fmt.Sprintf("error shutting down server: %s", err.Error()))
	}
}
