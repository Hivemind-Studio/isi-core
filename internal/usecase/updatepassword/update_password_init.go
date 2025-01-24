package updatepassword

type UseCase struct {
	repoUser repoUserInterface
}

func NewUpdatePasswordUseCase(
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoUser,
	}
}
