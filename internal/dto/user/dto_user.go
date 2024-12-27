package user

import (
	"time"
)

type UserDTO struct {
	Name        string     `json:"name,omitempty"`
	Email       string     `json:"email"`
	Address     *string    `json:"address,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string     `json:"gender"`
	Occupation  *string    `json:"occupation"`
	Status      *bool      `json:"status"`
	Role        *string    `json:"role"`
	Photo       *string    `json:"photo"`
	CreatedAt   *time.Time `json:"created_at"`
}

type GetUsersDTO struct {
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      *int64     `json:"createrole"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type SuspendDTO struct {
	Ids           []int64 `json:"ids"`
	UpdatedStatus string  `json:"updated_status"`
}
