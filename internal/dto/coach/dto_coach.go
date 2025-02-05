package coach

import "time"

type CreateCoachDTO struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	Address     string `json:"address" validate:"required"`
}

type PatchPasswordCoach struct {
	Password          string `json:"password"`
	ConfirmPassword   string `json:"confirm_password"`
	VerificationToken string `db:"verification_token"`
}

type DTO struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Address        *string    `json:"address"`
	PhoneNumber    *string    `json:"phone_number"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Gender         *string    `json:"gender"`
	Occupation     *string    `json:"occupation,omitempty"`
	Status         int64      `json:"status"`
	Role           *string    `json:"role,omitempty"`
	Photo          *string    `json:"photo"`
	CreatedAt      *time.Time `json:"created_at"`
	Certifications *string    `json:"certifications,omitempty"`
	Experiences    *string    `json:"experiences,omitempty"`
	Educations     *string    `json:"educations,omitempty"`
	Level          *string    `json:"level"`
	Title          *string    `json:"title"`
	Bio            *string    `json:"bio"`
	Expertise      *string    `json:"expertise"`
}

type QueryCoachDTO struct {
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Status      string     `json:"status"`
	Level       string     `json:"level"`
	Role        *int64     `json:"role"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}
