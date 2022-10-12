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
	auth := r.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/sign-up", h.singUpUser).Methods(http.MethodPost)
	auth.HandleFunc("/sign-in", h.signInUser).Methods(http.MethodPost)

	s := r.PathPrefix("/api").Subrouter()
	s.Use(h.validateAuthHeader)
	//chats
	s.HandleFunc("/chat", h.createChat).Methods(http.MethodPost)
	s.HandleFunc("/chat", h.getAllChats).Methods(http.MethodGet)
	//messages
	s.HandleFunc("/chat/{id:[0-9]+}/message", h.getAllChatMessages).Methods(http.MethodGet)
	//users
	s.HandleFunc("/user/{user}", h.findUserByName).Methods(http.MethodGet)
}
