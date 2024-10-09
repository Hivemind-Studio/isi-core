package user

import "github.com/go-playground/validator"

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ValidateLoginReq(loginReq LoginDTO) error {
	validate := validator.New()
	return validate.Struct(loginReq)
}

type RegisterDTO struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

func ValidateRegisterReq(registerDto RegisterDTO) error {
	validate := validator.New()
	return validate.Struct(registerDto)
}
