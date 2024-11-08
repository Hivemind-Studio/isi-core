package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"time"
)

type UserDTO struct {
	Name        string       `json:"name,omitempty"`
	Email       string       `json:"email"`
	Address     *string      `json:"address,omitempty"`
	PhoneNumber *string      `json:"phone_number,omitempty"`
	DateOfBirth *time.Time   `json:"date_of_birth"`
	Gender      *enum.Gender `json:"gender"`
	Occupation  *string      `json:"occupation"`
	Status      *bool        `json:"status"`
	Role        *string      `json:"role"`
	Photo       *string      `json:"photo"`
	CreatedAt   *time.Time   `json:"created_at"`
}

func ConvertUsersToDTOs(users []user.User) []UserDTO {
	dtos := make([]UserDTO, len(users))
	for i, u := range users {
		dtos[i] = UserDTO{
			Name:        u.Name,
			Email:       u.Email,
			Address:     u.Address,
			PhoneNumber: u.PhoneNumber,
			DateOfBirth: u.DateOfBirth,
			Gender:      u.Gender,
			Occupation:  u.Occupation,
			Status:      u.Status,
		}
	}
	return dtos
}

func ConvertUserToDTO(user user.User) UserDTO {
	return UserDTO{
		Name:        user.Name,
		Email:       user.Email,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		DateOfBirth: user.DateOfBirth,
		Gender:      user.Gender,
		Occupation:  user.Occupation,
		Status:      user.Status,
	}
}
