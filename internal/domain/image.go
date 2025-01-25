package domain

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID        uuid.UUID  `db:"id"`
	Image     []byte     `db:"image"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
