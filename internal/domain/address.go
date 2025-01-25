package domain

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID        uuid.UUID   `db:"id"`
	Name      string      `db:"name"`
	Street    string      `db:"street"`
	Type      AddressType `db:"type"`
	CityID    uuid.UUID   `db:"city_id"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
	DeletedAt *time.Time  `db:"deleted_at"`
	City      City        `db:"-"`
}

type AddressType int

const (
	AddressTypeClient   AddressType = 0
	AddressTypeSupplier AddressType = 10
)
