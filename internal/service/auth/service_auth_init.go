package auth

type AuthService struct {
	repoAuth repoAuthInterface
}

func NewAuthService(repoAuth repoAuthInterface) *AuthService {
	return &AuthService{
		repoAuth,
	}
}
