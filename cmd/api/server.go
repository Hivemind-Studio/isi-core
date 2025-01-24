package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	handleauth "github.com/Hivemind-Studio/isi-core/internal/handler/http/auth"
	handlecoach "github.com/Hivemind-Studio/isi-core/internal/handler/http/coach"
	handlecoachee "github.com/Hivemind-Studio/isi-core/internal/handler/http/coachee"
	handleprofile "github.com/Hivemind-Studio/isi-core/internal/handler/http/profile"
	handlerole "github.com/Hivemind-Studio/isi-core/internal/handler/http/role"
	handleuser "github.com/Hivemind-Studio/isi-core/internal/handler/http/user"
	repoCoach "github.com/Hivemind-Studio/isi-core/internal/repository/coach"
	reporole "github.com/Hivemind-Studio/isi-core/internal/repository/role"
	repouser "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/internal/service/useremail"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createcoach"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createrole"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createstaff"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createuser"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/forgotpassword"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getcoachbyid"
	getcoaceehbyid "github.com/Hivemind-Studio/isi-core/internal/usecase/getcoacheebyid"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getcoachees"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getcoaches"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getprofileuser"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getuserbyid"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getusers"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/sendverification"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updatepassword"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updateprofile"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updateprofilepassword"
	updateuserole "github.com/Hivemind-Studio/isi-core/internal/usecase/updateuserrole"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updateuserstatus"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/userlogin"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/verifyregistrationtoken"

	"github.com/gofiber/fiber/v2"
)

type AppApi struct {
	authHandle    *handleauth.Handler
	userHandle    *handleuser.Handler
	roleHandle    *handlerole.Handler
	coachHandle   *handlecoach.Handler
	coacheeHandle *handlecoachee.Handler
	profileHandle *handleprofile.Handler
}

type Router interface {
	RegisterRoutes(app *fiber.App)
}

func routerList(app *AppApi) []Router {
	return []Router{
		app.authHandle,
		app.coachHandle,
		app.userHandle,
		app.roleHandle,
		app.coacheeHandle,
		app.profileHandle,
	}
}

func initApp(cfg *configs.Config) (*AppApi, error) {
	dbConn := dbInitConnection(cfg)
	emailClient := initEmailClient(cfg)

	userRepo := repouser.NewUserRepo(dbConn)
	roleRepo := reporole.NewRoleRepo(dbConn)
	coachRepo := repoCoach.NewCoachRepo(dbConn)

	createCoachUseCase := createcoach.NewCreateCoachUseCase(coachRepo, userRepo, emailClient)
	getCoachByIdUseCase := getcoachbyid.NewGetCoachByIdUseCase(coachRepo)

	userEmailService := useremail.NewUserEmailService(userRepo, emailClient)

	createRoleUseCase := createrole.NewCreateRoleUseCase(roleRepo)
	userLoginUseCase := userlogin.NewLoginUseCase(userRepo)
	sendVerificationUseCase := sendverification.NewSendVerificationUseCase(userRepo, userEmailService)
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
	getProfileUser := getprofileuser.NewGetProfileUserByLogin(userRepo)
	updateProfilePassword := updateprofilepassword.NewUpdateProfilePasswordUseCase(userRepo)
	updateProfile := updateprofile.NewUpdateProfileUseCase(userRepo)

	roleHandler := handlerole.NewRoleHandler(createRoleUseCase)
	authHandler := handleauth.NewAuthHandler(userLoginUseCase,
		sendVerificationUseCase,
		verificationRegistrationTokenUseCase,
		createUserUseCase,
		updateCoachPasswordUseCase,
		forgotPasswordUseCase)
	userHandler := handleuser.NewUserHandler(
		createUserStaffUseCase,
		getUsersUseCase,
		getUserByIdUseCase,
		updateUserStatusUseCase,
		updateUserRoleUseCase)
	coachHandler := handlecoach.NewCoachHandler(getCoachesUseCase, createCoachUseCase, getCoachByIdUseCase)
	coacheeHandler := handlecoachee.NewCoacheeHandler(getCoacheesUseCase, getCoacheeByIdUseCase)
	profileHandler := handleprofile.NewProfileHandler(getProfileUser, updateProfilePassword, updateProfile)

	return &AppApi{
			userHandle:    userHandler,
			roleHandle:    roleHandler,
			authHandle:    authHandler,
			coacheeHandle: coacheeHandler,
			coachHandle:   coachHandler,
			profileHandle: profileHandler,
		},
		nil
}
