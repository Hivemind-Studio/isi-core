package verifyregistrationtoken

type UseCase struct {
	repoUser repoUserInterface
}

func NewVerifyRegistrationTokenUsecase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
