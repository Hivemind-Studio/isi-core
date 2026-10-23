package getuserbyid

type UseCase struct {
	repoUser repoUserInterface
}

func NewGetUserByIdUseCase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
