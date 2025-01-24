package updateprofile

type UseCase struct {
	repoUser repoUserInterface
}

func NewUpdateProfileUseCase(
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoUser,
	}
}
