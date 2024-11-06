package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"time"
)

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationDTO struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required,min=4"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=4"`
	Email           string `json:"email" validate:"required"`
}

type RegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Photo string `json:"photo"`
}

type UserDTO struct {
	ID          int64        `json:"id"`
	RoleId      int64        `json:"role_id"`
	Name        string       `json:"name,omitempty"`
	Email       string       `json:"email"`
	Address     *string      `json:"address,omitempty"`
	PhoneNumber *string      `json:"phone_number,omitempty"`
	DateOfBirth *time.Time   `json:"date_of_birth"`
	Gender      *enum.Gender `json:"gender"`
	Occupation  *string      `json:"occupation"`
	Status      *bool        `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
}

func ConvertUsersToDTOs(users []user.User) []UserDTO {
	dtos := make([]UserDTO, len(users))
	for i, user := range users {
		dtos[i] = UserDTO{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Address:     user.Address,
			PhoneNumber: user.PhoneNumber,
			DateOfBirth: user.DateOfBirth,
			Gender:      user.Gender,
			Occupation:  user.Occupation,
			Status:      user.Status,
			CreatedAt:   user.CreatedAt,
		}
	}
	return dtos
}
