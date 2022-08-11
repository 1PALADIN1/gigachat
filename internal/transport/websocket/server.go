package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
	http.HandleFunc("/", srv.handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", srv.port), nil); err != nil {
		return err
	}

	return nil
}

func (srv *WebSocketServer) handler(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in connection:", err)
		return
	}
	defer connection.Close()

	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		log.Println(string(message))

		greeting := fmt.Sprintf("Hello from server, %v!", connection.RemoteAddr())
		connection.WriteMessage(1, []byte(greeting))
	}
}
