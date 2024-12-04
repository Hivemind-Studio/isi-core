package auth

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationDTO struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required,min=4"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=4"`
	Email           string `json:"email" validate:"required"`
	Token           string `json:"token" validate:"required"`
}

type EmailVerificationDTO struct {
	Email string `json:"email" validate:"required"`
}

type CoachRegistrationDTO struct {
	Password        string `json:"password" validate:"required,min=4"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=4"`
	Token           string `json:"token" validate:"required"`
}

type RegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Photo string `json:"photo"`
	Token string `json:"token"`
}
