package repository

import (
	"context"

	"github.com/mestvl-shop-app/backend/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Client ClientRepositoryInterface
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Client: newClientRepository(db),
	}
}

type ClientRepositoryInterface interface {
	Create(ctx context.Context, client *domain.Client) error
}
