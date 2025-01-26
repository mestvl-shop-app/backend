package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	authv1 "github.com/mestvl-shop-app/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

const op = "auth service grpc client"

type Client struct {
	api authv1.AuthClient
	log *slog.Logger
}

type ClientInterface interface {
	Register(ctx context.Context, in *RegisterInput) (*uuid.UUID, error)
	Login(ctx context.Context, in *LoginInput) (string, error)
	Validate(ctx context.Context, token string) (bool, error)
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	// Setup logging
	rpcLogger := log.With("service", "gRPC/client")
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(interceptorLogger(rpcLogger), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	}

	conn, err := grpc.NewClient(addr, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient failed: %w", err)
	}

	return &Client{
		api: authv1.NewAuthClient(conn),
		log: log,
	}, nil
}

// interceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

type RegisterInput struct {
	Email    string
	Password string
}

func (c *Client) Register(ctx context.Context, in *RegisterInput) (*uuid.UUID, error) {
	resp, err := c.api.Register(ctx, &authv1.RegisterRequest{
		Email:    in.Email,
		Password: in.Password,
	})

	if err != nil {
		if errors.Is(err, status.Error(codes.AlreadyExists, "user already exists")) {
			return nil, ErrClientAlreadyExists
		}

		return nil, fmt.Errorf("%s: register new client failed: %w", op, err)
	}

	idStr := resp.GetUserId()
	id := uuid.MustParse(idStr)

	return &id, nil
}

type LoginInput struct {
	Email    string
	Password string
	AppID    int
}

func (c *Client) Login(ctx context.Context, in *LoginInput) (string, error) {
	resp, err := c.api.Login(ctx, &authv1.LoginRequest{
		Email:    in.Email,
		Password: in.Password,
		AppId:    int32(in.AppID),
	})
	if err != nil {
		if errors.Is(err, status.Error(codes.InvalidArgument, "invalid email or password")) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("%s: login failed: %w", op, err)
	}

	return resp.GetToken(), nil
}

func (c *Client) Validate(ctx context.Context, token string) (bool, error) {
	res, err := c.api.Validate(ctx, &authv1.ValidateRequest{Token: token})
	if err != nil {
		return false, fmt.Errorf("%s: validate failed: %w", op, err)
	}

	if res.GetStatus() == authv1.ValidateStatus_FORBIDDEN {
		return false, nil
	}

	return true, nil
}
