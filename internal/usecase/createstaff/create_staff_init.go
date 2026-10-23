package createstaff

type UseCase struct {
	repoUser         repoUserInterface
	userEmailService userEmailService
}

func NewCreateUserStaffUseCase(
	repoUser repoUserInterface,
	userEmailService userEmailService,
) *UseCase {
	return &UseCase{
		repoUser,
		userEmailService,
	}
}
