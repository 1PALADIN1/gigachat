package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1PALADIN1/gigachat_server/log/internal/repository"
	"github.com/1PALADIN1/gigachat_server/log/internal/service"
	"github.com/1PALADIN1/gigachat_server/log/internal/transport/srv_grpc"
)

type Config struct {
	Server struct {
		GRPCPort int
	}
}

func Run(config *Config) {
	repo := repository.NewRepository()
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
