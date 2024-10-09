package role

type Service struct {
	repoRole repoRoleInterface
}

func NewRoleService(repoRole repoRoleInterface) *Service {
	return &Service{
		repoRole: repoRole,
	}
}
