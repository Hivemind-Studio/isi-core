package updatecoachpassword

type UseCase struct {
	repoCoach repoCoachInterface
	repoUser  repoUserInterface
}

func NewUpdateCoachPasswordUseCase(
	repoCoach repoCoachInterface,
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoCoach,
		repoUser,
	}
}
