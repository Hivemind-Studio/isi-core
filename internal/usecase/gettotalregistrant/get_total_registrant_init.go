package gettotalregistrant

type UseCase struct {
	repoUser repoUserInterface
}

func NewGetTotalRegistrantUseCase(
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoUser,
	}
}
