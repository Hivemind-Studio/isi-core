package coach

type Service struct {
	repoUser repoUserInterface
}

func NewCoacheeService(repoUser repoUserInterface) *Service {
	return &Service{
		repoUser: repoUser,
	}
}
