package sendconfirmationchangenewemail

type UseCase struct {
	repoUser     repoUserInterface
	emailService userEmailService
}

func NewSendConfirmationChangeNewEmail(
	repoUser repoUserInterface,
	emailService userEmailService,
) *UseCase {
	return &UseCase{
		repoUser,
		emailService,
	}
}
