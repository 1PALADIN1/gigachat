package websocket

import (
	"fmt"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/transport/websocket/handlers"
)

func NewServer(port int, handlerAddress string) *http.Server {
	mux := http.NewServeMux()
	handler := handlers.Handler{}
	mux.HandleFunc(handlerAddress, handler.DefaultHandler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
