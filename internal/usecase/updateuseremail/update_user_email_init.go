package updateuseremail

type UseCase struct {
	repoUser repoUserInterface
}

func NewUpdateUserEmailUseCase(
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoUser,
	}
}
