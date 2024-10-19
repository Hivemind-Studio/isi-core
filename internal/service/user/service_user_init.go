package user

type UserService struct {
	repoUser repoUserInterface
}

func NewUserService(repoUser repoUserInterface) *UserService {
	return &UserService{
		repoUser: repoUser,
	}
}
