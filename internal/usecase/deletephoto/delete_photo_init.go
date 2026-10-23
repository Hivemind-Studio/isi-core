package deletephoto

type UseCase struct {
	repoUser repoUserInterface
}

func NewDeletePhotoStatusUseCase(repoUser repoUserInterface) *UseCase {
	return &UseCase{
		repoUser,
	}
}
