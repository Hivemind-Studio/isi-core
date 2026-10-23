package googleoauth2

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OauthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

func NewGoogleOauth(config *OauthConfig) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  config.RedirectURL,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
