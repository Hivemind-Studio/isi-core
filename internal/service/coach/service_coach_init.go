package coach

type Service struct {
	repoCoach repoCoachInterface
}

func NewCoachService(repoCoach repoCoachInterface) *Service {
	return &Service{
		repoCoach,
	}
}
