package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1PALADIN1/gigachat_server/log/internal/repository"
	"github.com/1PALADIN1/gigachat_server/log/internal/repository/mysql"
	"github.com/1PALADIN1/gigachat_server/log/internal/service"
	"github.com/1PALADIN1/gigachat_server/log/internal/transport/srv_grpc"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Server struct {
		GRPCPort int
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
	defer stopService(db)

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := srv_grpc.NewHandler(service)

	go func() {
		if err := handler.ListenGRPC(config.Server.GRPCPort); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("LogService started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func setupDB(config *Config) (*sqlx.DB, error) {
	db, err := mysql.NewDB(config.DB.DSN, float64(config.DB.ConnectionTimeout))
	if err != nil {
		return nil, err
	}

	return db, err
}

func stopService(db *sqlx.DB) {
	log.Println("LogService shutting down")

	if err := db.Close(); err != nil {
		log.Printf("error closing db: %s", err.Error())
	}
}
