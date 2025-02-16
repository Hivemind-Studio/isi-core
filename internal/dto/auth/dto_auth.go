package auth

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationDTO struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=12"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=12"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	Email           string `json:"email,omitempty"`
	Token           string `json:"token" validate:"required"`
	Gender          string `json:"gender,omitempty"`
	Address         string `json:"address,omitempty"`
}

type GoogleCallbackDTO struct {
	Code string `json:"code"`
}

type GoogleCallbackResponse struct {
	Token string `json:"token"`
}

type GoogleLoginResponse struct {
	Url string `json:"url"`
}

type EmailVerificationDTO struct {
	Email string `json:"email" validate:"required"`
}

type RegistrationStaffDTO struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone_number" validate:"required"`
	Email   string `json:"email" validate:"required"`
	Address string `json:"address" validate:"required"`
	Gender  string `json:"gender" validate:"required"`
	Role    string `json:"role" validate:"required"`
}

type UpdatePasswordRegistration struct {
	Password        string `json:"password" validate:"required,min=8,max=12"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=12"`
	Token           string `json:"token" validate:"required"`
}

type RegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Photo string `json:"photo"`
	Token string `json:"token"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdatePassword struct {
	CurrentPassword string `json:"current" validate:"required,min=8,max=12"`
	Password        string `json:"password" validate:"required,min=8,max=12"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=12"`
}
