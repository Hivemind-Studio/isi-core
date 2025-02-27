package deletecampaign

type UseCase struct {
	repoCampaign repoCampaignInterface
}

func NewDeleteCampaignUseCase(repoCampaign repoCampaignInterface) *UseCase {
	return &UseCase{
		repoCampaign,
	}
}
