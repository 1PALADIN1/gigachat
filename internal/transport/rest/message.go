package rest

import (
	"net/http"
	"strconv"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
	"github.com/gorilla/mux"
)

// Получение всех сообщений из чата
func (h *Handler) getAllChatMessages(w http.ResponseWriter, r *http.Request) {
	_, ok := helper.ValidateAuthHeader(w, r, h.service.Authorization)
	if !ok {
		return
	}

	chatId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, "invalid input param")
		return
	}

	messages, err := h.service.GetAllMessages(chatId)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]entity.ResponseMessage, 0)
	for _, m := range messages {
		resp = append(resp, m.ToResponse())
	}

	helper.SendResponse(w, http.StatusOK, resp)
}
