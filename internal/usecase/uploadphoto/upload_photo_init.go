package uploadphoto

type UseCase struct {
	repoUser repoUserInterface
}

func NewUpdatePhotoStatusUseCase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
