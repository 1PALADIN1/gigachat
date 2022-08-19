package helper

import (
	"net/http"
)

func ValidateRequestMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	isValid := r.Method == method
	if !isValid {
		SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}

	return isValid
}
