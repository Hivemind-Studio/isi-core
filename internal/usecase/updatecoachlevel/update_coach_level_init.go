package updatecoachlevel

type UseCase struct {
	repoCoach repoCoachInterface
}

func NewUpdateCoachLevelUseCase(repoCoach repoCoachInterface) *UseCase {
	return &UseCase{
		repoCoach,
	}
}
