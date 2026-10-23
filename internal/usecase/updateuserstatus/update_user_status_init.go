package updateuserstatus

type UseCase struct {
	repoUser repoUserInterface
}

func NewUpdateUserStatusUseCase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
