package googleoauthcallback

import "golang.org/x/oauth2"

type GoogleOAuthUseCase struct {
	OAuthConfig *oauth2.Config
	repoUser    repoUserInterface
	APIUrl      string
}

func NewGoogleOAuthUseCase(config *oauth2.Config, repoUser repoUserInterface, apiURL string,
) *GoogleOAuthUseCase {
	return &GoogleOAuthUseCase{
		repoUser:    repoUser,
		OAuthConfig: config,
		APIUrl:      apiURL,
	}
}
