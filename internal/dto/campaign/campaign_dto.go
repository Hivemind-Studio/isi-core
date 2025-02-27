package campaign

import "time"

type DTO struct {
	Name      string    `json:"name" validate:"required"`
	Channel   string    `json:"channel" validate:"required"`
	Link      string    `json:"link" validate:"required"`
	Status    bool      `json:"status" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}
