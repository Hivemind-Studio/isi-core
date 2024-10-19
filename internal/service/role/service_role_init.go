package role

type RoleService struct {
	repoRole repoRoleInterface
}

func NewRoleService(repoRole repoRoleInterface) *RoleService {
	return &RoleService{
		repoRole: repoRole,
	}
}
