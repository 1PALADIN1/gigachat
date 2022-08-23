package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
	"github.com/gorilla/mux"
)

// Получение всех сообщений из чата
func (h *Handler) getAllChatMessages(w http.ResponseWriter, r *http.Request) {
	userId, ok := helper.ValidateAuthHeader(w, r, h.service.Authorization)
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

	user, err := h.service.GetUserById(userId)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("user with id %d is not found", userId))
		return
	}

	resp := make([]entity.ResponseMessage, 0)
	for _, m := range messages {
		resp = append(resp, m.ToResponse(user))
	}

	helper.SendResponse(w, http.StatusOK, resp)
}
