package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	handleauth "github.com/Hivemind-Studio/isi-core/internal/handler/http/auth"
	handlecoach "github.com/Hivemind-Studio/isi-core/internal/handler/http/coach"
	handlecoachee "github.com/Hivemind-Studio/isi-core/internal/handler/http/coachee"
	handleprofile "github.com/Hivemind-Studio/isi-core/internal/handler/http/profile"
	handleuser "github.com/Hivemind-Studio/isi-core/internal/handler/http/user"
	repoCoach "github.com/Hivemind-Studio/isi-core/internal/repository/coach"
	repouser "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/internal/service/useremail"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createcoach"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createstaff"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createuser"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/deletephoto"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/forgotpassword"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getcoachbyid"
	getcoaceehbyid "github.com/Hivemind-Studio/isi-core/internal/usecase/getcoacheebyid"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getcoachees"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getcoaches"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getprofileuser"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getuserbyid"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getusers"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/googlelogin"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/googleoauthcallback"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/sendchangeemailverification"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/sendconfirmationchangenewemail"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/sendregistrationverification"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updatecoachlevel"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updatepassword"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updateprofile"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updateprofilepassword"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updateuseremail"
	updateuserole "github.com/Hivemind-Studio/isi-core/internal/usecase/updateuserrole"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updateuserstatus"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/uploadphoto"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/userlogin"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/verifyregistrationtoken"
	"github.com/Hivemind-Studio/isi-core/pkg/session"

	"github.com/gofiber/fiber/v2"
)

type AppApi struct {
	authHandle    *handleauth.Handler
	userHandle    *handleuser.Handler
	coachHandle   *handlecoach.Handler
	coacheeHandle *handlecoachee.Handler
	profileHandle *handleprofile.Handler
}

type Router interface {
	RegisterRoutes(app *fiber.App, sessionManager *session.SessionManager)
}

func routerList(app *AppApi) []Router {
	return []Router{
		app.authHandle,
		app.coachHandle,
		app.userHandle,
		app.coacheeHandle,
		app.profileHandle,
	}
}

func initApp(cfg *configs.Config, sessionManager *session.SessionManager) (*AppApi, error) {
	dbConn := dbInitConnection(cfg)
	emailClient := initEmailClient(cfg)
	googleOauthClient := initGoogleOauthClient(cfg)

	userRepo := repouser.NewUserRepo(dbConn)
	coachRepo := repoCoach.NewCoachRepo(dbConn)

	createCoachUseCase := createcoach.NewCreateCoachUseCase(coachRepo, userRepo, emailClient)
	getCoachByIdUseCase := getcoachbyid.NewGetCoachByIdUseCase(coachRepo)

	userEmailService := useremail.NewUserEmailService(userRepo, emailClient)

	userLoginUseCase := userlogin.NewLoginUseCase(userRepo)
	sendVerificationUseCase := sendregistrationverification.NewSendVerificationUseCase(userRepo, userEmailService)
	verificationRegistrationTokenUseCase := verifyregistrationtoken.NewVerifyRegistrationTokenUsecase(userRepo)
	createUserUseCase := createuser.NewCreateUserUseCase(userRepo)
	updateUserStatusUseCase := updateuserstatus.NewUpdateUserStatusUseCase(userRepo)
	updateCoachPasswordUseCase := updatepassword.NewUpdatePasswordUseCase(userRepo)
	getUsersUseCase := getusers.NewGetUsersUseCase(userRepo)
	getUserByIdUseCase := getuserbyid.NewGetUserByIdUseCase(userRepo)
	getCoachesUseCase := getcoaches.NewGetCoachesUseCase(coachRepo)
	getCoacheeByIdUseCase := getcoaceehbyid.NewGetCoacheeInterface(userRepo)

	getCoacheesUseCase := getcoachees.NewGetCoacheesUseCase(userRepo)
	createUserStaffUseCase := createstaff.NewCreateUserStaffUseCase(userRepo, userEmailService)

	forgotPasswordUseCase := forgotpassword.NewForgotPasswordUseCase(userRepo, userEmailService)
	updateUserRoleUseCase := updateuserole.NewUpdateUserRoleUseCase(userRepo)
	googleLoginUseCase := googlelogin.NewGoogleLoginUseCase(googleOauthClient)
	googleOAuthCallbackUseCase := googleoauthcallback.NewGoogleOAuthCallbackUseCase(googleOauthClient, userRepo)
	getProfileUser := getprofileuser.NewGetProfileUserByLogin(userRepo, coachRepo)
	updateProfilePassword := updateprofilepassword.NewUpdateProfilePasswordUseCase(userRepo)
	updateProfile := updateprofile.NewUpdateProfileUseCase(userRepo, coachRepo)
	updateUserEmail := updateuseremail.NewUpdateUserEmailUseCase(userRepo)
	sendChangeEmailVerification := sendchangeemailverification.NewSendChangeEmailVerificationUseCase(userRepo, userEmailService)
	sendConfirmationChangeNewEmail := sendconfirmationchangenewemail.NewSendConfirmationChangeNewEmail(userRepo, userEmailService)
	uploadPhoto := uploadphoto.NewUpdatePhotoStatusUseCase(userRepo)
	deletePhoto := deletephoto.NewDeletePhotoStatusUseCase(userRepo)
	updateCoachLevel := updatecoachlevel.NewUpdateCoachLevelUseCase(coachRepo)

	authHandler := handleauth.NewAuthHandler(
		sessionManager,
		userLoginUseCase,
		sendVerificationUseCase,
		verificationRegistrationTokenUseCase,
		createUserUseCase,
		updateCoachPasswordUseCase,
		forgotPasswordUseCase,
		googleLoginUseCase,
		googleOAuthCallbackUseCase,
	)
	userHandler := handleuser.NewUserHandler(
		createUserStaffUseCase,
		getUsersUseCase,
		getUserByIdUseCase,
		updateUserStatusUseCase,
		updateUserRoleUseCase,
		updateUserEmail,
		sendChangeEmailVerification,
		sendConfirmationChangeNewEmail)
	coachHandler := handlecoach.NewCoachHandler(getCoachesUseCase, createCoachUseCase, getCoachByIdUseCase, updateCoachLevel)
	coacheeHandler := handlecoachee.NewCoacheeHandler(getCoacheesUseCase, getCoacheeByIdUseCase)
	profileHandler := handleprofile.NewProfileHandler(getProfileUser, updateProfilePassword, updateProfile,
		uploadPhoto, deletePhoto)

	return &AppApi{
			authHandle:    authHandler,
			userHandle:    userHandler,
			coacheeHandle: coacheeHandler,
			coachHandle:   coachHandler,
			profileHandle: profileHandler,
		},
		nil
}
