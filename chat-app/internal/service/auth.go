package service

import (
	"context"
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService struct {
	addr        string
	connTimeout time.Duration
}

func NewAuthService(addr string, connTimeout time.Duration) *AuthService {
	return &AuthService{
		addr:        addr,
		connTimeout: connTimeout,
	}
}

// SignUpUser регистрирует нового пользователя.
// Возвращает ID ползователя в БД в случае успеха. В противном случае - ошибку.
func (s *AuthService) SignUpUser(user entity.User) (int, error) {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), s.connTimeout)
	defer cancel()

	response, err := c.SignUpUser(ctx, &auth.SignUpRequest{
		User: &auth.User{
			Username: user.Username,
			Password: user.Password,
		},
	})

	if err != nil {
		return 0, err
	}

	return int(response.Id), nil
}

// GenerateToken генерирует JWT-токен после успешной авторизации пользователя
func (s *AuthService) GenerateToken(username, password string) (string, int, error) {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return "", 0, err
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), s.connTimeout)
	defer cancel()

	response, err := c.GenerateToken(ctx, &auth.GenerateTokenRequest{
		User: &auth.User{
			Username: username,
			Password: password,
		},
	})

	if err != nil {
		return "", 0, err
	}

	return response.Token, int(response.Id), nil
}

// ParseToken валидирует токен и возвращает id пользователя в случае успеха
func (s *AuthService) ParseToken(token string) (int, error) {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), s.connTimeout)
	defer cancel()

	response, err := c.ParseToken(ctx, &auth.ParseTokenRequest{
		Token: token,
	})

	if err != nil {
		return 0, err
	}

	return int(response.Id), nil
}
