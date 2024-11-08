package auth

type Service struct {
	repoAuth repoUserInterface
}

func NewAuthService(repoAuth repoUserInterface) *Service {
	return &Service{
		repoAuth,
	}
}
