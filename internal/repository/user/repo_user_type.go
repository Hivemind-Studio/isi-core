package user

import "time"

type User struct {
	ID          int64     `db:"id" json:"id"`
	Username    string    `db:"username" json:"username" validate:"required,min=3,max=50"`
	Password    string    `db:"password" json:"-"`
	Name        string    `db:"name" json:"name,omitempty"`
	Email       string    `db:"email" json:"email" validate:"required,email"`
	Address     string    `db:"address" json:"address,omitempty"`
	PhoneNumber string    `db:"phone_number" json:"phone_number,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy   int64     `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy   int64     `db:"updated_by" json:"updated_by,omitempty"`
}
