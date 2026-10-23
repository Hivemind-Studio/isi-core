package forgotpassword

type UseCase struct {
	repoUser         repoUserInterface
	userEmailService userEmailService
}

func NewForgotPasswordUseCase(repoUser repoUserInterface, userEmailService userEmailService,
) *UseCase {
	return &UseCase{
		repoUser,
		userEmailService,
	}
}
