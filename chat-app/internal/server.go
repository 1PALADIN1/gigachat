package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/logger"
)

type Server struct {
	httpServer *http.Server
}

type ServerConfig struct {
	Port           int
	Handler        http.Handler
	ReadTimeout    int
	WriteTimeout   int
	MaxHeaderBytes int
}

func NewServer(config ServerConfig) *Server {
	logger.LogInfo(fmt.Sprintf("Setup server at port %d. Read timeout: %d [sec], write timout: %d [sec]", config.Port, config.ReadTimeout, config.WriteTimeout))

	server := new(Server)
	server.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Port),
		Handler:        config.Handler,
		MaxHeaderBytes: config.MaxHeaderBytes,
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
	}

	return server
}

func (s *Server) Start() error {
	logger.LogInfo("Starting server...")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.LogInfo("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
