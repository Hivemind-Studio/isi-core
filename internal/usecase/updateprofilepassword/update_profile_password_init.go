package updateprofilepassword

type UseCase struct {
	repoUser repoUserInterface
}

func NewUpdateProfilePasswordUseCase(
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoUser,
	}
}
