package getcoaches

type UseCase struct {
	repoCoach repoCoachInterface
}

func NewGetCoachesUseCase(repoCoach repoCoachInterface) *UseCase {
	return &UseCase{
		repoCoach,
	}
}
