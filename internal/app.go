// Стартовая точка приложения.
// Здесь создаются основные сущности, строятся зависимости.
package internal

import (
	"log"

	"github.com/1PALADIN1/gigachat_server/internal/transport/websocket"
)

const (
	defaultPort = 8000 //TODO: вынести в конфиг
)

var (
	server Server
)

type Server interface {
	Start() error
}

func Run() {
	server = websocket.NewServer(defaultPort)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
