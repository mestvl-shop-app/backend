package domain

import (
	"time"

	"github.com/google/uuid"
)

type SupplierAddress struct {
	ID         uuid.UUID  `db:"id"`
	SupplierID uuid.UUID  `db:"supplier_id"`
	AddressID  uuid.UUID  `db:"address_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
	Supplier   Supplier   `db:"-"`
	Address    Address    `db:"-"`
}
