package srv_grpc

import (
	"context"

	"github.com/1PALADIN1/gigachat_server/log/internal/service"
	"github.com/1PALADIN1/gigachat_server/log/internal/transport/srv_grpc/logs"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	service service.Log
}

// Log логирует сообщения
func (s *LogServer) Log(ctx context.Context, rq *logs.LogRequest) (*logs.LogResponse, error) {
	if err := s.service.Log(rq.LogLevel, rq.Message); err != nil {
		return nil, err
	}

	return &logs.LogResponse{}, nil
}
