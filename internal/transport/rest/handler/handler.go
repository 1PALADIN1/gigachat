package handler

import (
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth", authUser)
}
