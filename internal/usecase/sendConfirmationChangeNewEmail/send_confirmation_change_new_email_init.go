package updateuseremail

type UseCase struct {
	repoUser     repoUserInterface
	emailService userEmailService
}

func NewUpdateUserEmailUseCase(
	repoUser repoUserInterface,
	emailService userEmailService,
) *UseCase {
	return &UseCase{
		repoUser,
		emailService,
	}
}
