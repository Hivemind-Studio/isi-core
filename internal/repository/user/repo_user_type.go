package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"time"
)

type User struct {
	ID           int64        `db:"id"`
	RoleId       *int64       `db:"role_id"`
	Password     string       `db:"password"`
	Name         string       `db:"name"`
	Email        string       `db:"email"`
	Address      *string      `db:"address"`
	PhoneNumber  *string      `db:"phone_number"`
	DateOfBirth  *time.Time   `db:"date_of_birth"`
	Gender       *enum.Gender `db:"gender"`
	Verification *bool        `db:"verification"`
	Occupation   *string      `db:"occupation"`
	Photo        *string      `db:"photo"`
	Status       *bool        `db:"status"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}

type Cookie struct {
	ID       int64   `db:"id"`
	Password string  `db:"password"`
	Name     string  `db:"name"`
	RoleId   *int64  `db:"role_id"`
	RoleName *string `db:"role_name"`
	Email    *string `db:"email"`
	Photo    *string `db:"photo"`
}
