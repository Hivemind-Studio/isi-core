package getcoachbyid

type UseCase struct {
	repoCoach repoCoachInterface
}

func NewGetCoachByIdUseCase(repoCoach repoCoachInterface) *UseCase {
	return &UseCase{
		repoCoach,
	}
}
