package websocket

import (
	"log"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/service"
	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
	"github.com/gorilla/websocket"
)

type Handler struct {
	service *service.Service
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws", h.setupWsConnection)
}

func (h *Handler) setupWsConnection(w http.ResponseWriter, r *http.Request) {
	userId, ok := helper.ValidateAuthHeader(w, r, h.service.Authorization)
	if !ok {
		return
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in connection:", err)
		return
	}

	go func() {
		defer connection.Close()
		h.handleUserMessages(connection, userId)
	}()
}
