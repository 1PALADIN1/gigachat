// Стартовая точка приложения.
// Здесь создаются основные сущности, строятся зависимости.
package internal

import (
	"log"
	"net/http"
	"sync"

	"github.com/1PALADIN1/gigachat_server/internal/service"

	"github.com/1PALADIN1/gigachat_server/internal/transport/rest"
	"github.com/1PALADIN1/gigachat_server/internal/transport/websocket"
)

type ServerInfo struct {
	server *http.Server
	name   string
}

type Config struct {
	Server struct {
		Ws struct {
			Port    int
			Handler string
		}
		Rest struct {
			Port int
		}
	}
	Token struct {
		SigningKey string `yaml:"signing-key"`
	}
}

func Run(config *Config) {
	service.Init(config.Token.SigningKey)

	wg := new(sync.WaitGroup)
	serverConfig := config.Server

	servers := []ServerInfo{
		ServerInfo{
			websocket.NewServer(serverConfig.Ws.Port, serverConfig.Ws.Handler),
			"WebSocket",
		},
		ServerInfo{
			rest.NewServer(serverConfig.Rest.Port),
			"REST",
		},
	}

	for _, server := range servers {
		wg.Add(1)
		go startServer(server.name, server.server, wg)
	}

	wg.Wait()
}

func startServer(name string, server *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Starting %s server at addr: %v\n", name, server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
