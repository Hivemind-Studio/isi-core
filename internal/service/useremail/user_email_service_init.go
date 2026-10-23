package useremail

import "github.com/Hivemind-Studio/isi-core/pkg/mail"

type Service struct {
	repoUser    repoUserInterface
	emailClient *mail.EmailClient
}

func NewUserEmailService(repoUser repoUserInterface, emailClient *mail.EmailClient,
) *Service {
	return &Service{
		repoUser,
		emailClient,
	}
}
