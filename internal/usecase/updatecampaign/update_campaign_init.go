package updatecampaign

type UseCase struct {
	repoCampaign repoCampaignInterface
}

func NewUpdateCampaignUseCase(
	repoCampaign repoCampaignInterface,
) *UseCase {
	return &UseCase{
		repoCampaign,
	}
}
