package campaign

import "time"

type Campaign struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Channel   string    `db:"channel"`
	StartDate string    `db:"start_date"`
	EndDate   string    `db:"end_date"`
	Link      string    `db:"link"`
	Status    bool      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
