package internal

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1PALADIN1/gigachat_server/auth/internal/repository"
	"github.com/1PALADIN1/gigachat_server/auth/internal/repository/postgres"
	"github.com/1PALADIN1/gigachat_server/auth/internal/service"
	"github.com/1PALADIN1/gigachat_server/auth/internal/transport/srv_grpc"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Server struct {
		GRPCPort int
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

	service := service.NewService(repo, srvConfig)
	handler := srv_grpc.NewHandler(service)

	go func() {
		if err := handler.ListenGRPC(config.Server.GRPCPort); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("AuthService started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	stopService(db)
}

func setupDB(config *Config) (*sqlx.DB, error) {
	db, err := postgres.NewDB(config.DB.DSN, float64(config.DB.ConnectionTimeout))
	if err != nil {
		return nil, err
	}

	return db, err
}

func stopService(db *sqlx.DB) {
	log.Println("AuthService shutting down")

	if err := db.Close(); err != nil {
		log.Printf("error closing db: %s", err.Error())
	}
}
