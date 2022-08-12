package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/service"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketServer struct {
	port int
}

func NewServer(port int) *WebSocketServer {
	server := WebSocketServer{
		port: port,
	}

	return &server
}

func (srv *WebSocketServer) Start() error {
	log.Println("Starting WebSocket server")

	http.HandleFunc("/", srv.mainHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", srv.port), nil); err != nil {
		return err
	}

	return nil
}

func (srv *WebSocketServer) mainHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in connection:", err)
		return
	}

	log.Println("Client", connection.RemoteAddr(), "connected!")
	service.StartUserSession(connection)
}
