package auth

type LoginDTO struct {
	Email    string `json:"useremail"`
	Password string `json:"password"`
}

type RegistrationDTO struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required,min=4"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=4"`
	Email           string `json:"useremail" validate:"required"`
	Token           string `json:"token" validate:"required"`
}

type EmailVerificationDTO struct {
	Email string `json:"useremail" validate:"required"`
}

type RegistrationStaffDTO struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone_number" validate:"required"`
	Email   string `json:"useremail" validate:"required"`
	Address string `json:"address" validate:"required"`
	Gender  string `json:"gender" validate:"required"`
	Role    string `json:"role" validate:"required"`
}

type CoachRegistrationDTO struct {
	Password        string `json:"password" validate:"required,min=4"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=4"`
	Token           string `json:"token" validate:"required"`
}

type RegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"useremail"`
}

type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"useremail"`
	Role  string `json:"createrole"`
	Photo string `json:"photo"`
	Token string `json:"token"`
}
