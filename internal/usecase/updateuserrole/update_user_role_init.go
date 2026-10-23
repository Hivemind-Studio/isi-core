package updateuserole

type UseCase struct {
	repoUser repoUserInterface
}

func NewUpdateUserRoleUseCase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
