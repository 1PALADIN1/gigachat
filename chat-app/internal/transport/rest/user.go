package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
	"github.com/gorilla/mux"
)

// Поиск пользователей по username
func (h *Handler) findUserByName(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("error parsing user id: %s", err.Error()))
		return
	}

	filter, ok := mux.Vars(r)["user"]
	if !ok {
		helper.SendErrorResponse(w, http.StatusBadRequest, "invalid input param")
		return
	}

	users, err := h.service.FindUserByName(filter, userId)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	respUsers := make([]entity.ResponseUser, 0)
	for _, u := range users {
		respUsers = append(respUsers, u.ToResponse())
	}

	helper.SendResponse(w, http.StatusOK, respUsers)
}
