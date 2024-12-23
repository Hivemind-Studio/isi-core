package createstaff

import "github.com/Hivemind-Studio/isi-core/pkg/mail"

type UseCase struct {
	repoUser         repoUserInterface
	userEmailService userEmailService
	emailClient      *mail.EmailClient
}

func NewCreateUserStaffUseCase(repoUser repoUserInterface, userEmailService userEmailService, emailClient *mail.EmailClient) *UseCase {
	return &UseCase{
		repoUser,
		userEmailService,
		emailClient,
	}
}
