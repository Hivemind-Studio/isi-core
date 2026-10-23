package createusercampaign

type UseCase struct {
	repoCampaign repoCampaignInterface
	repoUser     repoUserInterface
}

func NewCreateUserCampaignUseCase(
	repoCampaign repoCampaignInterface,
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoCampaign,
		repoUser,
	}
}
