package coach

import (
	"github.com/Hivemind-Studio/isi-core/pkg/mail"
)

type Service struct {
	repoCoach   repoCoachInterface
	repoUser    repoUserInterface
	emailClient *mail.EmailClient
}

func NewCoachService(repoCoach repoCoachInterface, repoUser repoUserInterface, emailClient *mail.EmailClient) *Service {
	return &Service{
		repoCoach,
		repoUser,
		emailClient,
	}
}
