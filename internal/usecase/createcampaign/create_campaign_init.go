package createcampaign

type UseCase struct {
	repoCampaign repoCampaignInterface
	repoUser     repoUserInterface
}

func NewCreateCampaignUseCase(
	repoCampaign repoCampaignInterface,
	repoUser repoUserInterface,
) *UseCase {
	return &UseCase{
		repoCampaign,
		repoUser,
	}
}
