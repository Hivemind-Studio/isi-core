package getcoachees

type UseCase struct {
	repoUser repoUserInterface
}

func NewGetCoacheesUseCase(
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoUser,
	}
}
