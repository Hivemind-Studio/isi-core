package user

import (
	"time"
)

type UserDTO struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Address        *string    `json:"address"`
	PhoneNumber    *string    `json:"phone_number"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Gender         *string    `json:"gender"`
	Occupation     *string    `json:"occupation"`
	Status         int64      `json:"status"`
	Role           *string    `json:"role,omitempty"`
	Photo          *string    `json:"photo"`
	CreatedAt      *time.Time `json:"created_at"`
	Title          *string    `json:"title,omitempty"`
	Bio            *string    `json:"bio,omitempty"`
	Expertise      *string    `json:"expertise,omitempty"`
	Certifications *string    `json:"certifications,omitempty"`
	Experiences    *string    `json:"experiences,omitempty"`
	Educations     *string    `json:"educations,omitempty"`
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

type UserRole struct {
	Role string `json:"role"`
}

type ChangeEmailDTO struct {
	Token    string `json:"token"`
	NewEmail string `json:"new_email"`
}

type ConfirmChangeEmailDTO struct {
	Token string `json:"token"`
}

type ChangeEmailEmailVerificationDTO struct {
	Email string `json:"email" validate:"required"`
}

type UpdateUserDTO struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	Occupation  string `json:"occupation,omitempty"`
	Title       string `json:"title,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Expertise   string `json:"expertise,omitempty"`
}

type UpdateCoachDTO struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	Title       string `json:"title"`
	Bio         string `json:"bio"`
	Expertise   string `json:"expertise"`
}
