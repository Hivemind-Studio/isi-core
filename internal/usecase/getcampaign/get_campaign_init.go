package getcampaign

type UseCase struct {
	repoCampaign repoCampaignInterface
}

func NewGetCampaignUseCase(
	repoCampaign repoCampaignInterface,
) *UseCase {
	return &UseCase{
		repoCampaign,
	}
}
