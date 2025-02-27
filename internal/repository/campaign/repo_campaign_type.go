package campaign

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"time"
)

type Campaign struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	Channel   string     `db:"channel"`
	StartDate *time.Time `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`
	Link      string     `db:"link"`
	Status    int8       `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type UserRegistration struct {
	ID               int64     `db:"id"`
	UserID           int64     `db:"user_id"`
	CampaignID       int64     `db:"campaign_id"`
	RegistrationDate time.Time `db:"registration_date"`
	UserAgent        string    `db:"user_agent"`
	IPAddress        string    `db:"ip_address"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func ConvertCampaignToDTO(c Campaign) campaign.DTO {
	return campaign.DTO{
		ID:        c.ID,
		Name:      c.Name,
		Channel:   c.Channel,
		Status:    c.Status,
		Link:      c.Link,
		StartDate: c.StartDate,
		EndDate:   c.EndDate,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func ConvertCampaignToDTOs(c []Campaign) []campaign.DTO {
	dtos := make([]campaign.DTO, len(c))
	for i, u := range c {
		dtos[i] = ConvertCampaignToDTO(u)
	}
	return dtos
}
