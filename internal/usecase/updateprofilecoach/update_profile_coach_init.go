package updateprofilecoach

type UseCase struct {
	repoCoach repoCoachInterface
}

func NewUpdateProfileCoachUseCase(
	repoCoach repoCoachInterface,
) *UseCase {
	return &UseCase{
		repoCoach,
	}
}
