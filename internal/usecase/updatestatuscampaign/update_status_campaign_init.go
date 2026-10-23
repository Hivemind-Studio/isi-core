package updatestatuscampaign

type UseCase struct {
	repoCampaign repoCampaignInterface
}

func NewUpdateStatusCampaignUseCase(repoCampaign repoCampaignInterface) *UseCase {
	return &UseCase{
		repoCampaign,
	}
}
