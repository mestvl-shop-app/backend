package service

import (
	"context"
	"log/slog"

	"github.com/mestvl-shop-app/backend/internal/client/auth"
	"github.com/mestvl-shop-app/backend/internal/config"
	"github.com/mestvl-shop-app/backend/internal/repository"
)

type Services struct {
	Client ClientServiceInterface
}

type Deps struct {
	Logger     *slog.Logger
	Config     *config.Config
	Repos      *repository.Repositories
	AuthClient *auth.Client
}

func NewServices(deps Deps) *Services {
	return &Services{
		Client: newClientService(
			deps.Repos.Client,
			deps.AuthClient,
			deps.Config,
		),
	}
}

type ClientServiceInterface interface {
	Register(ctx context.Context, input *RegisterClientInput) error
	Login(ctx context.Context, email string, password string) (string, error)
}
