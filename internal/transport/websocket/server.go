package websocket

import (
	"fmt"
	"net/http"
)

type WebSocketServer struct {
	port int
}

func NewServer(port int) *WebSocketServer {
	server := WebSocketServer{
		port: port,
	}

	return &server
}

func (w *WebSocketServer) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world! Path: %s", r.URL.Path)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", w.port), nil); err != nil {
		return err
	}

	return nil
}
