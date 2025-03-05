package getcampaignbyid

type UseCase struct {
	repoCampaign repoCampaignInterface
}

func NewGetCampaignByIdUseCase(repoCampaign repoCampaignInterface) *UseCase {
	return &UseCase{
		repoCampaign,
	}
}
