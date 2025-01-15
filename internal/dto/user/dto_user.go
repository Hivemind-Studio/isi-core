package user

import (
	"time"
)

type UserDTO struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Address     *string    `json:"address"`
	PhoneNumber *string    `json:"phone_number"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      *string    `json:"gender"`
	Occupation  *string    `json:"occupation"`
	Status      bool       `json:"status"`
	Role        *string    `json:"role,omitempty"`
	Photo       *string    `json:"photo"`
	CreatedAt   *time.Time `json:"created_at"`
}

type GetUsersDTO struct {
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Status      string     `json:"status"`
	Level       string     `json:"level"`
	Role        *int64     `json:"role"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type SuspendDTO struct {
	Ids           []int64 `json:"ids"`
	UpdatedStatus string  `json:"updated_status"`
}
