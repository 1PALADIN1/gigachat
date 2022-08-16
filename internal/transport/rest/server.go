package rest

import (
	"fmt"
	"net/http"

	"github.com/1PALADIN1/gigachat_server/internal/transport/rest/handler"
)

func NewServer(port int) *http.Server {
	mux := http.NewServeMux()
	handler.SetupRoutes(mux)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
