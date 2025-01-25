package repository

import (
	"context"
	"fmt"

	"github.com/mestvl-shop-app/backend/internal/db"
	"github.com/mestvl-shop-app/backend/internal/domain"

	"github.com/jmoiron/sqlx"
)

type clientRepository struct {
	db *sqlx.DB
}

func newClientRepository(db *sqlx.DB) *clientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client *domain.Client) error {
	const query = "" +
		"INSERT INTO client " +
		"(id, firstname, surname, birthday, gender) " +
		"VALUES($1, $2, $3, $4, $5);"

	_, err := r.db.ExecContext(ctx, query, client.ID, client.Firstname, client.Surname, client.Birthday, client.Gender)

	if err != nil {
		if db.IsDuplicate(err) {
			return domain.ErrDuplicateEntry
		}
		return fmt.Errorf("insert client failed: %w", err)
	}

	return nil
}
