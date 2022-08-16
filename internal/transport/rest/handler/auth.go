package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/service"
)

type tokenResponse struct {
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
}

func authUser(w http.ResponseWriter, r *http.Request) {
	resp := &tokenResponse{}

	token, err := service.GenerateToken()
	if err != nil {
		log.Println("error generating token:", err)
		resp.Error = "Token is not generated"
		sendTokenResponse(resp, w, http.StatusUnauthorized)
		return
	}

	resp.AccessToken = token
	sendTokenResponse(resp, w, http.StatusOK)
}

func sendTokenResponse(resp *tokenResponse, w http.ResponseWriter, statusCode int) {
	jsonString, err := json.Marshal(resp)
	if err != nil {
		log.Println("error marshaling json:", err)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(jsonString)
	if err != nil {
		log.Println("error writing response:", err)
	}
}
