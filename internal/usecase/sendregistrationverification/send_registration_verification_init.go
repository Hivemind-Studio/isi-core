package sendregistrationverification

type UseCase struct {
	repoUser         repoUserInterface
	userEmailService userEmailService
}

func NewSendVerificationUseCase(
	repoUser repoUserInterface,
	userEmailService userEmailService,
) *UseCase {
	return &UseCase{
		repoUser,
		userEmailService,
	}
}
