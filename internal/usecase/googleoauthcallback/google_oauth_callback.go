package googleoauthcallback

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"log"
	"net/http"
)

func (uc *GoogleOAuthUseCase) GetUserDataFromGoogle(ctx context.Context, code string) (user.User, error) {
	token, err := uc.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		return user.User{}, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get(uc.APIUrl + token.AccessToken)
	if err != nil {
		return user.User{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	var userInfo struct {
		ID     string `json:"id"`
		Email  string `json:"email"`
		Name   string `json:"name"`
		Avatar string `json:"picture"`
	}

	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		log.Println("Failed to parse user info")
		return user.User{}, httperror.New(http.StatusInternalServerError, "Failed to parse user info")
	}

	log.Printf("Succeed fetch User info: %+v", userInfo)

	users, err := uc.repoUser.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		return user.User{}, err
	}

	if *users.GoogleID == userInfo.ID {
		return user.User{}, nil
	}

	return user.User{}, nil
}
