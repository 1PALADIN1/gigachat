package rest

import (
	"encoding/json"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/entity"

	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
)

// Хендлер регистрации нового пользователя
func (h *Handler) singUpUser(w http.ResponseWriter, r *http.Request) {
	if !helper.ValidateRequestMethod(w, r, http.MethodPost) {
		return
	}
	defer r.Body.Close()

	var input entity.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		helper.SendErrorResponse(w, http.StatusBadRequest, err.Error())
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
