package updateprofile

type UseCase struct {
	repoUser  repoUserInterface
	repoCoach repoCoachInterface
}

func NewUpdateProfileUseCase(
	repoUser repoUserInterface,
	repoCoach repoCoachInterface,
) *UseCase {
	return &UseCase{
		repoUser,
		repoCoach,
	}
}
