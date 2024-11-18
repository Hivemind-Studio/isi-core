package coach

type Service struct {
	repoUser repoUserInterface
}

func NewCoachService(repoUser repoUserInterface) *Service {
	return &Service{
		repoUser: repoUser,
	}
}
