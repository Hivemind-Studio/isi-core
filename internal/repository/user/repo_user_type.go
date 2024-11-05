package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"time"
)

type User struct {
	ID           int64        `db:"id" json:"id"`
	Password     string       `db:"password" json:"-"`
	Name         string       `db:"name" json:"name,omitempty"`
	Email        string       `db:"email" json:"email" validate:"required,email"`
	Address      *string      `db:"address" json:"address,omitempty"`
	PhoneNumber  *string      `db:"phone_number" json:"phone_number,omitempty"`
	DateOfBirth  *time.Time   `db:"date_of_birth" json:"date_of_birth"`
	Gender       *enum.Gender `db:"gender" json:"gender"`
	Verification *bool        `db:"verification" json:"verification"`
	Occupation   *string      `db:"occupation" json:"occupation"`
	Status       *bool        `db:"status" json:"status"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
}
