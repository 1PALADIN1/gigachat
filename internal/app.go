// Стартовая точка приложения.
// Здесь создаются основные сущности, строятся зависимости.
package internal

import (
	"log"

	"github.com/1PALADIN1/gigachat_server/internal/transport/websocket"
)

type Server interface {
	Start() error
}

type Config struct {
	Server struct {
		Port      int
		WsAddress string `yaml:"ws-address"`
	}
}

var (
	server Server
)

func Run(config *Config) {
	server = websocket.NewServer(config.Server.Port, config.Server.WsAddress)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
