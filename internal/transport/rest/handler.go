package rest

import (
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SetupRoutes(r *mux.Router) {
	//auth
	r.HandleFunc("/api/auth/sign-up", h.singUpUser).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/sign-in", h.signInUser).Methods(http.MethodPost)
	//chats
	r.HandleFunc("/api/chat", h.createChat).Methods(http.MethodPost)
	r.HandleFunc("/api/chat", h.getAllChats).Methods(http.MethodGet)
	//messages
	r.HandleFunc("/api/chat/{id:[0-9]+}", h.getAllChatMessages).Methods(http.MethodGet)
	//users
	r.HandleFunc("/api/user/{user}", h.findUserByName).Methods(http.MethodGet)
}
