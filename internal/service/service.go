package service

import (
	"context"
	"log/slog"

	"github.com/mestvl-shop-app/backend/internal/config"
	"github.com/mestvl-shop-app/backend/internal/repository"
	"github.com/mestvl-shop-app/backend/pkg/auth"
	hash "github.com/mestvl-shop-app/backend/pkg/hasher"
)

type Services struct {
	Client ClientServiceInterface
}

type Deps struct {
	Logger       *slog.Logger
	Config       *config.Config
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
	Repos        *repository.Repositories
}

func NewServices(deps Deps) *Services {
	return &Services{
		Client: newClientService(deps.Repos.Client),
	}
}

type ClientServiceInterface interface {
	Register(ctx context.Context, dto *RegisterClientDTO) error
}
