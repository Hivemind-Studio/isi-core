package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	userdto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type User struct {
	ID             int64      `db:"id"`
	RoleId         *int64     `db:"role_id"`
	Password       *string    `db:"password"`
	Name           string     `db:"name"`
	Email          string     `db:"email"`
	GoogleID       *string    `db:"google_id"`
	Address        *string    `db:"address"`
	PhoneNumber    *string    `db:"phone_number"`
	DateOfBirth    *time.Time `db:"date_of_birth"`
	Gender         *string    `db:"gender"`
	Verification   *bool      `db:"verification"`
	Occupation     *string    `db:"occupation"`
	Photo          *string    `db:"photo"`
	Status         int64      `db:"status"`
	Version        int64      `db:"version"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	Certifications *string    `db:"certifications"`
	Experiences    *string    `db:"experiences"`
	Educations     *string    `db:"educations"`
	Title          *string    `db:"title"`
	Bio            *string    `db:"bio"`
	Expertise      *string    `db:"expertise"`
	Level          *int64     `db:"level"`
	RoleName       *string    `db:"role_name"`
}

type EmailVerification struct {
	Id                int64     `db:"id"`
	Email             string    `db:"email"`
	VerificationToken string    `db:"verification_token"`
	Trial             int8      `db:"trial"`
	Version           int64     `db:"version"`
	ExpiredAt         time.Time `db:"expired_at"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func ConvertUserToDTO(user User) userdto.UserDTO {
	roleName := constant.GetRoleName(user.RoleId)
	return userdto.UserDTO{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		DateOfBirth: user.DateOfBirth,
		Gender:      user.Gender,
		Occupation:  user.Occupation,
		Status:      user.Status,
		Photo:       user.Photo,
		Role:        &roleName,
		RoleID:      user.RoleId,
		Title:       user.Title,
		Bio:         user.Bio,
		Expertise:   user.Expertise,
		CreatedAt:   user.CreatedAt,
	}
}

func ConvertUsersToDTOs(users []User) []userdto.UserDTO {
	dtos := make([]userdto.UserDTO, len(users))
	for i, u := range users {
		dtos[i] = ConvertUserToDTO(u)
	}
	return dtos
}
