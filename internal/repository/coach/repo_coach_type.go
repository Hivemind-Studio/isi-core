package coach

type PatchPassword struct {
	Password          string `json:"password"`
	VerificationToken string `db:"verification_token"`
}
