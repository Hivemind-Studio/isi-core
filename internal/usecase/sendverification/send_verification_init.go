package sendverification

import "github.com/Hivemind-Studio/isi-core/pkg/mail"

type UseCase struct {
	repoUser    repoUserInterface
	emailClient *mail.EmailClient
}

func NewSendVerificationUseCase(repoUser repoUserInterface, emailClient *mail.EmailClient) *UseCase {
	return &UseCase{
		repoUser,
		emailClient,
	}
}
