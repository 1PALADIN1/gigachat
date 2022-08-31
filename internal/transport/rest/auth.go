package rest

import (
	"encoding/json"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/entity"

	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
)

type successSignInResponse struct {
	UserId int    `json:"id"`
	Token  string `json:"access_token"`
}

// Хендлер регистрации нового пользователя
func (h *Handler) singUpUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input entity.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.SignUpUser(input)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendResponse(w, http.StatusOK,
		map[string]interface{}{
			"id": id,
		})
}

// Хендлер авторизации пользователя
func (h *Handler) signInUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input entity.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, userId, err := h.service.GenerateToken(input.Username, input.Password)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendResponse(w, http.StatusOK, successSignInResponse{
		UserId: userId,
		Token:  token,
	})
}
