package srv_grpc

import (
	"context"

	"github.com/1PALADIN1/gigachat_server/auth/internal/entity"
	"github.com/1PALADIN1/gigachat_server/auth/internal/service"
	"github.com/1PALADIN1/gigachat_server/auth/internal/transport/srv_grpc/auth"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	service service.Authorization
}

// SignUpUser регистрирует нового пользователя
func (s *AuthServer) SignUpUser(ctx context.Context, rq *auth.SignUpRequest) (*auth.SignUpResponse, error) {
	input := rq.GetUser()
	user := entity.User{
		Username: input.Username,
		Password: input.Password,
	}

	userId, err := s.service.SignUpUser(user)
	if err != nil {
		rs := &auth.SignUpResponse{
			Id: 0,
		}

		return rs, err
	}

	rs := &auth.SignUpResponse{
		Id: int32(userId),
	}
	return rs, nil
}

// GenerateToken авторизирует пользователя и генерирует токен
func (s *AuthServer) GenerateToken(ctx context.Context, rq *auth.GenerateTokenRequest) (*auth.GenerateTokenResponse, error) {
	input := rq.GetUser()
	token, userId, err := s.service.GenerateToken(input.Username, input.Password)
	if err != nil {
		rs := &auth.GenerateTokenResponse{
			Token: "",
			Id:    0,
		}

		return rs, err
	}

	rs := &auth.GenerateTokenResponse{
		Token: token,
		Id:    int32(userId),
	}
	return rs, nil
}

// ParseToken валидирует токен и возварщает id пользователя в случае успеха
func (s *AuthServer) ParseToken(ctx context.Context, rq *auth.ParseTokenRequest) (*auth.ParseTokenResponse, error) {
	userId, err := s.service.ParseToken(rq.Token)
	if err != nil {
		rs := &auth.ParseTokenResponse{
			Id: 0,
		}

		return rs, err
	}

	rs := &auth.ParseTokenResponse{
		Id: int32(userId),
	}
	return rs, nil
}
