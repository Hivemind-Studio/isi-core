package createrole

type UseCase struct {
	repoRole repoRoleInterface
}

func NewCreateRoleUseCase(repoRole repoRoleInterface) *UseCase {
	return &UseCase{
		repoRole: repoRole,
	}
}
