package coach

import (
	coachDTO "github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"time"
)

type Coach struct {
	ID             int64      `db:"id"`
	Name           string     `db:"name"`
	Email          string     `db:"email"`
	Address        *string    `db:"address"`
	PhoneNumber    *string    `db:"phone_number"`
	DateOfBirth    *time.Time `db:"date_of_birth"`
	Gender         *string    `db:"gender"`
	Occupation     *string    `db:"occupation"`
	Photo          *string    `db:"photo"`
	Status         bool       `db:"status"`
	Verification   string     `db:"verification"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	Certifications string     `db:"certifications"`
	Experiences    string     `db:"experiences"`
	Educations     string     `db:"educations"`
	Level          string     `db:"level"`
	RoleID         int64      `db:"role_id"`
	RoleName       string     `db:"role_name"`
}

func ConvertCoachToDTO(coach Coach) coachDTO.DTO {
	return coachDTO.DTO{
		ID:             coach.ID,
		Name:           coach.Name,
		Email:          coach.Email,
		Address:        coach.Address,
		PhoneNumber:    coach.PhoneNumber,
		DateOfBirth:    coach.DateOfBirth,
		Gender:         coach.Gender,
		Occupation:     coach.Occupation,
		Status:         coach.Status,
		Photo:          coach.Photo,
		Certifications: coach.Certifications,
		Experiences:    coach.Experiences,
		Educations:     coach.Educations,
		Level:          coach.Level,
		CreatedAt:      &coach.CreatedAt,
	}
}

func ConvertCoachesToDTO(coaches []Coach) []coachDTO.DTO {
	dtos := make([]coachDTO.DTO, len(coaches))
	for i, u := range coaches {
		dtos[i] = ConvertCoachToDTO(u)
	}
	return dtos
}
