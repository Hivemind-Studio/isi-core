package googlelogin

import "golang.org/x/oauth2"

type UseCase struct {
	Oauth2 *oauth2.Config
}

func NewGoogleLoginUseCase(oauth2 *oauth2.Config) *UseCase {
	return &UseCase{
		Oauth2: oauth2,
	}
}
