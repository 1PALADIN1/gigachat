package helper

import (
	"net/http"
	"strings"

	"github.com/1PALADIN1/gigachat_server/internal/service"
)

const (
	authorizationHeader = "Authorization"
)

func ValidateRequestMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	isValid := r.Method == method
	if !isValid {
		SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}

	return isValid
}

func ValidateAuthHeader(w http.ResponseWriter, r *http.Request, authService service.Authorization) (int, bool) {
	header := r.Header.Get(authorizationHeader)
	if strings.TrimSpace(header) == "" {
		SendErrorResponse(w, http.StatusUnauthorized, "auth header is not set")
		return 0, false
	}

	// Authorization: Bearer <token>
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		SendErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
		return 0, false
	}

	userId, err := authService.ParseToken(headerParts[1])
	if err != nil {
		SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return 0, false
	}

	return userId, true
}
