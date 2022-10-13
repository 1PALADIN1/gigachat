package helper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/logger"
)

type errorResponse struct {
	Message string `json:"message"`
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logger.LogInfo(fmt.Sprintf("sending error: %s", message))
	SendResponse(w, statusCode, errorResponse{message})
}

func SendResponse(w http.ResponseWriter, statusCode int, message any) {
	resp, err := json.Marshal(message)
	if err != nil {
		logger.LogError(fmt.Sprintf("error marshaling response: %s", err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}
