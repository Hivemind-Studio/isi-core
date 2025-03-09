package googleoauthcallback

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/constant/loglevel"
	userDTO "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	user "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"net/http"
)

func (uc *GoogleOAuthCallbackUseCase) GetUserDataFromGoogle(ctx context.Context, code string) (userDTO.UserDTO, error) {
	googleToken, err := uc.OAuthConfig.Exchange(ctx, code)
	requestId := ctx.Value("request_id").(string)
	if err != nil {
		logger.Print(loglevel.ERROR, requestId, "google_oauth_callback",
			"GetUserDataFromGoogle", fmt.Sprintf("code exchange failed: %s", err.Error()), code)

		return userDTO.UserDTO{}, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	const GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	response, err := http.Get(GoogleUserInfoURL + googleToken.AccessToken)
	if err != nil {
		logger.Print(loglevel.ERROR, requestId, "google_oauth_callback",
			"GetUserDataFromGoogle", fmt.Sprintf("failed getting user info: %s", err.Error()),
			googleToken.AccessToken)
		return userDTO.UserDTO{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	var userInfo struct {
		ID     string `json:"id"`
		Email  string `json:"email"`
		Name   string `json:"name"`
		Avatar string `json:"picture"`
	}

	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		logger.Print(loglevel.ERROR, requestId, "google_oauth_callback",
			"GetUserDataFromGoogle", fmt.Sprintf("Failed to parse user info: %s", err.Error()), response.Body)
		return userDTO.UserDTO{}, httperror.New(http.StatusInternalServerError, "Failed to parse user info")
	}

	logger.Print(loglevel.INFO, requestId, "google_oauth_callback", "GetUserDataFromGoogle",
		fmt.Sprintf("Succeed fetch User info: %+v", userInfo), nil)

	savedUser, err := uc.repoUser.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		logger.Print(loglevel.ERROR, requestId, "google_oauth_callback",
			"FindByEmail", fmt.Sprintf("Error finding user by email: %s", err.Error()), userInfo.Email)

		if err.Error() == "user not found" {
			logger.Print(loglevel.INFO, requestId, "google_oauth_callback",
				"UserCreation", "User not found, proceeding with user creation", userInfo)

			tx, err := uc.repoUser.StartTx(ctx)
			if err != nil {
				logger.Print(loglevel.ERROR, requestId, "google_oauth_callback",
					"UserCreation", fmt.Sprintf("Failed to start transaction: %s", err.Error()), userInfo)
				return userDTO.UserDTO{}, err
			}

			createdUserId, err := uc.repoUser.Create(ctx, tx, user.CreateUserParams{
				Name:          userInfo.Name,
				Email:         userInfo.Email,
				Password:      nil,
				RoleID:        constant.RoleIDCoachee,
				PhoneNumber:   nil,
				Gender:        "",
				Address:       "",
				Status:        int(constant.ACTIVE),
				GoogleID:      &userInfo.ID,
				Photo:         &userInfo.Avatar,
				VerifiedEmail: true,
			})

			if err != nil {
				logger.Print(loglevel.ERROR, requestId, "google_oauth_callback",
					"UserCreation", fmt.Sprintf("Error creating user: %s", err.Error()), nil)
				dbtx.HandleRollback(tx)
				return userDTO.UserDTO{}, err
			}

			err = uc.repoUser.CommitTx()
			if err != nil {
				logger.Print(loglevel.ERROR, requestId, "google_oauth_callback",
					"UserCreation", fmt.Sprintf("Error committing transaction: %s", err.Error()), nil)
				return userDTO.UserDTO{}, err
			}

			logger.Print(loglevel.INFO, requestId, "google_oauth_callback",
				"UserCreation", fmt.Sprintf("User successfully created with ID: %d", createdUserId), createdUserId)

			savedUser, err = uc.repoUser.FindByEmail(ctx, userInfo.Email)
			if err != nil {
				return userDTO.UserDTO{}, err
			}

			return user.ConvertUserToDTO(savedUser), nil
		}

		return userDTO.UserDTO{}, err
	}

	if savedUser.GoogleID != nil && *savedUser.GoogleID == userInfo.ID {
		return user.ConvertUserToDTO(savedUser), nil
	}

	tx, err := uc.repoUser.StartTx(ctx)
	defer dbtx.HandleRollback(tx)
	if err != nil {
		return userDTO.UserDTO{}, err
	}

	err = uc.repoUser.UpdateUserGoogleId(ctx, tx, userInfo.Email, userInfo.ID)
	if err != nil {
		return userDTO.UserDTO{}, err
	}

	err = uc.repoUser.CommitTx()
	if err != nil {
		dbtx.HandleRollback(tx)
		return userDTO.UserDTO{}, err
	}

	return user.ConvertUserToDTO(savedUser), nil
}
