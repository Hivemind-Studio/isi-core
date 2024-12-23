package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	handleauth "github.com/Hivemind-Studio/isi-core/internal/handler/http/auth"
	handlecoach "github.com/Hivemind-Studio/isi-core/internal/handler/http/coach"
	handlecoachee "github.com/Hivemind-Studio/isi-core/internal/handler/http/coachee"
	handlerole "github.com/Hivemind-Studio/isi-core/internal/handler/http/role"
	handleuser "github.com/Hivemind-Studio/isi-core/internal/handler/http/user"
	repoCoach "github.com/Hivemind-Studio/isi-core/internal/repository/coach"
	reporole "github.com/Hivemind-Studio/isi-core/internal/repository/role"
	repouser "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createcoach"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createrole"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createstaff"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/createuser"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getcoachees"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getcoaches"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getuserbyid"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/getusers"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/sendverification"
	"github.com/Hivemind-Studio/isi-core/internal/usecase/updatecoachpassword"
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
}

type Router interface {
	RegisterRoutes(app *fiber.App)
}

func routerList(app *AppApi) []Router {
	return []Router{
		app.userHandle,
		app.authHandle,
		app.coachHandle,
		app.roleHandle,
		app.coacheeHandle,
	}
}

func initApp(cfg *configs.Config) (*AppApi, error) {
	dbConn := dbInitConnection(cfg)
	emailClient := initEmailClient(cfg)

	userRepo := repouser.NewUserRepo(dbConn)
	roleRepo := reporole.NewRoleRepo(dbConn)
	coachRepo := repoCoach.NewCoachRepo(dbConn)

	createRoleUseCase := createrole.NewCreateRoleUseCase(roleRepo)
	userLoginUseCase := userlogin.NewLoginUseCase(userRepo)
	sendVerificationUseCase := sendverification.NewSendVerificationUseCase(userRepo, emailClient)
	verificationRegistrationTokenUseCase := verifyregistrationtoken.NewVerifyRegistrationTokenUsecase(userRepo)
	createUserUseCase := createuser.NewCreateUserUseCase(userRepo)
	updateUserStatusUseCase := updateuserstatus.NewUpdateUserStatusUseCase(userRepo)
	updateCoachPasswordUseCase := updatecoachpassword.NewUpdateCoachPasswordUseCase(coachRepo, userRepo)
	getUsersUseCase := getusers.NewGetUsersUseCase(userRepo)
	getUserByIdUseCase := getuserbyid.NewGetUserByIdUseCase(userRepo)
	getCoachesUseCase := getcoaches.NewGetCoachesUseCase(coachRepo)
	createCoachUseCase := createcoach.NewCreateCoachUseCase(coachRepo, userRepo, emailClient)
	getCoacheesUseCase := getcoachees.NewGetCoacheesUseCase(userRepo)
	createUserStaffUseCase := createstaff.NewCreateUserStaffUseCase(userRepo, emailClient)

	roleHandler := handlerole.NewRoleHandler(createRoleUseCase)
	authHandler := handleauth.NewAuthHandler(userLoginUseCase,
		sendVerificationUseCase,
		verificationRegistrationTokenUseCase,
		createUserUseCase,
		updateCoachPasswordUseCase)
	userHandler := handleuser.NewUserHandler(
		createUserStaffUseCase,
		getUsersUseCase,
		getUserByIdUseCase,
		updateUserStatusUseCase)
	coachHandler := handlecoach.NewCoachHandler(getCoachesUseCase, createCoachUseCase)
	coacheeHandler := handlecoachee.NewCoacheeHandler(getCoacheesUseCase)

	return &AppApi{
			userHandle:    userHandler,
			roleHandle:    roleHandler,
			authHandle:    authHandler,
			coacheeHandle: coacheeHandler,
			coachHandle:   coachHandler,
		},
		nil
}
