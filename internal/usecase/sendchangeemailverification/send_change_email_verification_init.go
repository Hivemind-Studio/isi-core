package sendchangeemailverification

type UseCase struct {
	repoUser         repoUserInterface
	userEmailService userEmailService
}

func NewSendChangeEmailVerificationUseCase(
	repoUser repoUserInterface,
	userEmailService userEmailService,
) *UseCase {
	return &UseCase{
		repoUser,
		userEmailService,
	}
}
