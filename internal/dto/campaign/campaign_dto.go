package campaign

import "time"

type DTO struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name" validate:"required"`
	Channel       string     `json:"channel" validate:"required"`
	Link          string     `json:"link" validate:"required"`
	GeneratedLink string     `json:"generated_link,omitempty"`
	Status        int8       `json:"status" validate:"required"`
	StartDate     *time.Time `json:"start_date" validate:"required"`
	EndDate       *time.Time `json:"end_date" validate:"required"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type Params struct {
	Name      string     `json:"name" validate:"required"`
	Channel   string     `json:"channel" validate:"required"`
	Link      string     `json:"link" validate:"required"`
	Status    string     `json:"status" validate:"required"`
	StartDate *time.Time `json:"start_date" validate:"required"`
	EndDate   *time.Time `json:"end_date" validate:"required"`
}
