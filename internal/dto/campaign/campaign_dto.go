package campaign

import "time"

type DTO struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name" validate:"required"`
	Channel          string     `json:"channel" validate:"required"`
	Link             string     `json:"link" validate:"required"`
	GeneratedLink    string     `json:"generated_link,omitempty"`
	CampaignId       string     `json:"campaign_id" validate:"required"`
	Status           int8       `json:"status" validate:"required"`
	StartDate        *time.Time `json:"start_date" validate:"required"`
	EndDate          *time.Time `json:"end_date" validate:"required"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	TotalRegistrants *int64     `json:"total_registrants,omitempty"`
}

type Params struct {
	Name      string     `json:"name" validate:"required"`
	Channel   string     `json:"channel" validate:"required"`
	Link      string     `json:"link" validate:"required"`
	Status    string     `json:"status" validate:"required"`
	StartDate *time.Time `json:"start_date" validate:"required"`
	EndDate   *time.Time `json:"end_date" validate:"required"`
}

type PatchStatus struct {
	IDS    []int64 `json:"ids" validate:"required"`
	Status int8    `json:"status" validate:"required"`
}

type UserCampaign struct {
	Email      string `json:"email"`
	CampaignId string `json:"campaign_id"`
	UserAgent  string `json:"user_agent"`
	IPAddress  string `json:"ip_address"`
}
