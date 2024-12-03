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
	serviceauth "github.com/Hivemind-Studio/isi-core/internal/service/auth"
	servicecoach "github.com/Hivemind-Studio/isi-core/internal/service/coach"
	servicecoachee "github.com/Hivemind-Studio/isi-core/internal/service/coachee"
	servicerole "github.com/Hivemind-Studio/isi-core/internal/service/role"
	serviceuser "github.com/Hivemind-Studio/isi-core/internal/service/user"

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
		app.coachHandle,
		app.authHandle,
		app.userHandle,
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

	roleService := servicerole.NewRoleService(roleRepo)
	userService := serviceuser.NewUserService(userRepo)
	authService := serviceauth.NewAuthService(userRepo, emailClient)
	coachService := servicecoach.NewCoachService(coachRepo, emailClient)
	coacheeService := servicecoachee.NewCoacheeService(userRepo)

	roleHandler := handlerole.NewRoleHandler(roleService)
	authHandler := handleauth.NewAuthHandler(authService, userService)
	userHandler := handleuser.NewUserHandler(userService)
	coachHandler := handlecoach.NewCoachHandler(coachService)
	coacheeHandler := handlecoachee.NewCoacheeHandler(coacheeService)

	return &AppApi{
			userHandle:    userHandler,
			roleHandle:    roleHandler,
			authHandle:    authHandler,
			coacheeHandle: coacheeHandler,
			coachHandle:   coachHandler,
		},
		nil
}
