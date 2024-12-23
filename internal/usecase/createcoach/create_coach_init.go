package createcoach

import "github.com/Hivemind-Studio/isi-core/pkg/mail"

type UseCase struct {
	repoCoach   repoCoachInterface
	repoUser    repoUserInterface
	emailClient *mail.EmailClient
}

func NewCreateCoachUseCase(
	repoCoach repoCoachInterface,
	repoUser repoUserInterface,
	emailClient *mail.EmailClient,
) *UseCase {
	return &UseCase{
		repoCoach,
		repoUser,
		emailClient,
	}
}
