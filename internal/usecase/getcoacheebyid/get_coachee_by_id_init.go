package getcoaceehbyid

type UseCase struct {
	repoUser repoUserInterface
}

func NewGetCoacheeInterface(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
