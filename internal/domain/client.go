package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID        uuid.UUID     `db:"id"`
	Firstname string        `db:"firstname"`
	Surname   string        `db:"surname"`
	Birthday  *time.Time    `db:"birthday"`
	Gender    *ClientGender `db:"gender"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	DeletedAt *time.Time    `db:"deleted_at"`
	Addresses []Address     `db:"-"`
}

type ClientGender int

const (
	ClientGenderMale   ClientGender = 0
	ClientGenderFemale ClientGender = 10
)

type ClientGenderString string

const (
	ClientGenderStringMale   ClientGenderString = "male"
	ClientGenderStringFemale ClientGenderString = "female"
)

func (gender ClientGenderString) Code() ClientGender {
	switch gender {
	case ClientGenderStringMale:
		return ClientGenderMale
	case ClientGenderStringFemale:
		return ClientGenderFemale
	}

	return -1
}

func (gender *ClientGenderString) CodeFromPointer() *ClientGender {
	if gender == nil {
		return nil
	}

	var result ClientGender

	switch *gender {
	case ClientGenderStringMale:
		result = ClientGenderMale
		return &result
	case ClientGenderStringFemale:
		result = ClientGenderFemale
		return &result
	}

	return nil
}

func (gender ClientGender) String() ClientGenderString {
	switch gender {
	case ClientGenderMale:
		return ClientGenderStringMale
	case ClientGenderFemale:
		return ClientGenderStringFemale
	}

	return ""
}
