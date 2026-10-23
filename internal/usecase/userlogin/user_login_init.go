package userlogin

type UseCase struct {
	repoUser repoUserInterface
}

func NewLoginUseCase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
