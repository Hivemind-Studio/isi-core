package user

import (
	userdto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"time"
)

type User struct {
	ID           int64        `db:"id"`
	RoleId       *int64       `db:"role_id"`
	RoleName     *string      `db:"role_name"`
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

func ConvertUserToDTO(user User) userdto.UserDTO {
	return userdto.UserDTO{
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

func ConvertUsersToDTOs(users []User) []userdto.UserDTO {
	dtos := make([]userdto.UserDTO, len(users))
	for i, u := range users {
		dtos[i] = userdto.UserDTO{
			Name:        u.Name,
			Email:       u.Email,
			Address:     u.Address,
			PhoneNumber: u.PhoneNumber,
			DateOfBirth: u.DateOfBirth,
			Gender:      u.Gender,
			Occupation:  u.Occupation,
			Status:      u.Status,
			CreatedAt:   &u.CreatedAt,
		}
	}
	return dtos
}
