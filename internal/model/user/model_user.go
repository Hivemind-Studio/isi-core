package user

import "github.com/Hivemind-Studio/isi-core/internal/model/base"

type User struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"-" binding:"required"` // Hide password in response
	base.Model
}
