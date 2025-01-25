package domain

import (
	"time"

	"github.com/google/uuid"
)

type Supplier struct {
	ID          uuid.UUID  `db:"id"`
	Login       string     `db:"login"`
	Password    string     `db:"password"`
	Name        string     `db:"name"`
	PhoneNumber int        `db:"phone_number"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Addresses   []Address  `db:"-"`
}
