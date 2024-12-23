package coach

type CreateCoachDTO struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"useremail" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	Address     string `json:"address" validate:"required"`
}

type PatchPasswordCoach struct {
	Password          string `json:"password"`
	ConfirmPassword   string `json:"confirm_password"`
	VerificationToken string `db:"verification_token"`
}
