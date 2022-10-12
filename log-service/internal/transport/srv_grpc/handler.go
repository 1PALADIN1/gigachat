package srv_grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/1PALADIN1/gigachat_server/log/internal/service"
	"github.com/1PALADIN1/gigachat_server/log/internal/transport/srv_grpc/logs"
	"google.golang.org/grpc"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ListenGRPC прослушивает входящие сообщения по протоколу gRPC
func (h *Handler) ListenGRPC(portNumber int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		return fmt.Errorf("error starting listen grpc connections: %s", err.Error())
	}

	s := grpc.NewServer()
	logServer := &LogServer{
		service: h.service.Log,
	}
	logs.RegisterLogServiceServer(s, logServer)

	log.Println("gRPC server started at port", portNumber)

	if err := s.Serve(listener); err != nil {
		return fmt.Errorf("error serving grpc server: %s", err.Error())
	}

	return nil
}
