package websocket

import (
	"fmt"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/transport/websocket/handler"
)

func NewServer(port int, handlerAddress string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(handlerAddress, handler.DefaultHandler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
