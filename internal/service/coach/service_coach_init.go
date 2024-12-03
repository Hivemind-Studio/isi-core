package coach

import "github.com/Hivemind-Studio/isi-core/pkg/mail"

type Service struct {
	repoCoach   repoCoachInterface
	emailClient *mail.EmailClient
}

func NewCoachService(repoCoach repoCoachInterface, emailClient *mail.EmailClient) *Service {
	return &Service{
		repoCoach,
		emailClient,
	}
}
