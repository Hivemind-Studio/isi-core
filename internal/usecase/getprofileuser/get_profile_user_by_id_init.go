package getprofileuser

type UseCase struct {
	repoUser  repoUserInterface
	repoCoach repoCoachInterface
}

func NewGetProfileUserByLogin(repoUser repoUserInterface, repoCoach repoCoachInterface) *UseCase {
	return &UseCase{
		repoUser,
		repoCoach,
	}
}
