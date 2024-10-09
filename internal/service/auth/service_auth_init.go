package auth

type Service struct {
	repoAuth repoAuthInterface
}

func NewAuthService(repoAuth repoAuthInterface) *Service {
	return &Service{
		repoAuth,
	}
}
