package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	log.Println("sending error:", message)
	SendResponse(w, statusCode, errorResponse{message})
}

func SendResponse(w http.ResponseWriter, statusCode int, message any) {
	resp, err := json.Marshal(message)
	if err != nil {
		log.Println("error marshaling response:", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}
