package rest

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
	"github.com/gorilla/mux"
)

const (
	authorizationHeader = "Authorization"
)

// validateAuthHeader валидация заголовка входящего запроса
func (h *Handler) validateAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userId, err := h.parseHeader(r)
		if err != nil {
			helper.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		mux.Vars(r)["user_id"] = strconv.Itoa(userId)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) parseHeader(r *http.Request) (int, error) {
	header := r.Header.Get(authorizationHeader)
	if strings.TrimSpace(header) == "" {
		return 0, errors.New("invalid auth header")
	}

	// Authorization: Bearer <token>
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return 0, errors.New("invalid auth header")
	}

	if headerParts[0] != "Bearer" {
		return 0, errors.New("invalid auth header")
	}

	if strings.TrimSpace(headerParts[1]) == "" {
		return 0, errors.New("invalid auth header")
	}

	userId, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		return 0, err
	}

	return userId, err
}
