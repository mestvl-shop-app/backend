package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID               uuid.UUID  `db:"id"`
	Name             string     `db:"name"`
	CategoryID       uuid.UUID  `db:"category_id"`
	Price            float64    `db:"price"`
	AvailableStock   int        `db:"available_stock"`
	LastPurchaseDate *time.Time `db:"last_purchase_date"`
	SupplierID       uuid.UUID  `db:"supplier_id"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	Category         Category   `db:"-"`
	Supplier         Supplier   `db:"-"`
}
