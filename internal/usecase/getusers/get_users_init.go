package getusers

type UseCase struct {
	repoUser repoUserInterface
}

func NewGetUsersUseCase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
