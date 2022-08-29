package rest

import (
	"encoding/json"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
)

// Создание чата (или получение существующего, если такой чат уже существует)
func (h *Handler) createChat(w http.ResponseWriter, r *http.Request) {
	_, ok := helper.ValidateAuthHeader(w, r, h.service.Authorization)
	if !ok {
		return
	}
	defer r.Body.Close()

	var chat entity.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := chat.Validate(); err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	chatId, err := h.service.GetOrCreateChat(chat)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendResponse(w, http.StatusOK, map[string]interface{}{
		"id": chatId,
	})
}

// Получение всех чатов пользователя
func (h *Handler) getAllChats(w http.ResponseWriter, r *http.Request) {
	userId, ok := helper.ValidateAuthHeader(w, r, h.service.Authorization)
	if !ok {
		return
	}

	chats, err := h.service.GetAllChats(userId)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendResponse(w, http.StatusOK, chats)
}
