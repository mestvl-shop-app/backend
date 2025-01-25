package domain

import (
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	CountryID uuid.UUID  `db:"country_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Country   Country    `db:"-"`
}
