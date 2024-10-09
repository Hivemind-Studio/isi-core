package user

type Service struct {
	repoUser repoUserInterface
}

func NewUserService(repoUser repoUserInterface) *Service {
	return &Service{
		repoUser: repoUser,
	}
}
