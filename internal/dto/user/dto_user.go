package user

type LoginDTO struct {
	Email    string `json:"Email"`
	Password string `json:"password"`
}

type RegisterDTO struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type RegisterResponse struct {
	Name  string `json:"Name"`
	Email string `json:"email"`
}
