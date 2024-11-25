package auth

import "github.com/Hivemind-Studio/isi-core/pkg/mail"

type Service struct {
	repoUser    repoUserInterface
	emailClient *mail.EmailClient
}

func NewAuthService(repoUser repoUserInterface, emailClient *mail.EmailClient) *Service {
	return &Service{
		repoUser,
		emailClient,
	}
}
