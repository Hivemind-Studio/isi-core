package googleoauthcallback

import "golang.org/x/oauth2"

type GoogleOAuthCallbackUseCase struct {
	OAuthConfig *oauth2.Config
	repoUser    repoUserInterface
}

func NewGoogleOAuthCallbackUseCase(config *oauth2.Config, repoUser repoUserInterface,
) *GoogleOAuthCallbackUseCase {
	return &GoogleOAuthCallbackUseCase{
		repoUser:    repoUser,
		OAuthConfig: config,
	}
}
