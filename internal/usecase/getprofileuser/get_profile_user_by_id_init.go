package getprofileuser

type UseCase struct {
	repoUser repoUserInterface
}

func NewGetProfileUserByLogin(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
