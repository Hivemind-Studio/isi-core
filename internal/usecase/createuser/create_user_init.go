package createuser

type UseCase struct {
	repoUser repoUserInterface
}

func NewCreateUserUseCase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
