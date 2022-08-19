// Стартовая точка приложения.
// Здесь создаются основные сущности, строятся зависимости.
package internal

import (
	"log"
	"net/http"

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
		SigningKey string
	}
}

func Run(config Config) {
	mux := http.NewServeMux()
	wsHandler := websocket.NewHandler()
	wsHandler.SetupRoutes(mux)

	restHandler := rest.NewHandler()
	restHandler.SetupRoutes(mux)

	server := NewServer(ServerConfig{
		Port:           config.Server.Port,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, //1MB
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
	})

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
