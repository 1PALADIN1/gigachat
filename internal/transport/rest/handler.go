package rest

import (
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	//auth
	mux.HandleFunc("/api/auth/sign-up", h.singUpUser) //POST
	mux.HandleFunc("/api/auth/sign-in", h.signInUser) //POST
	//chats
	mux.HandleFunc("/api/chat", h.handleChat) //POST | GET
}
