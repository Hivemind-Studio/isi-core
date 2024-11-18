package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	handlecoach "github.com/Hivemind-Studio/isi-core/internal/handler/http/coach"
	handlecoachee "github.com/Hivemind-Studio/isi-core/internal/handler/http/coachee"
	handleuser "github.com/Hivemind-Studio/isi-core/internal/handler/http/user"
	repouser "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	servicecoach "github.com/Hivemind-Studio/isi-core/internal/service/coach"
	servicecoachee "github.com/Hivemind-Studio/isi-core/internal/service/coachee"
	serviceuser "github.com/Hivemind-Studio/isi-core/internal/service/user"

	handlerole "github.com/Hivemind-Studio/isi-core/internal/handler/http/role"
	reporole "github.com/Hivemind-Studio/isi-core/internal/repository/role"
	servicerole "github.com/Hivemind-Studio/isi-core/internal/service/role"

	handleauth "github.com/Hivemind-Studio/isi-core/internal/handler/http/auth"
	serviceauth "github.com/Hivemind-Studio/isi-core/internal/service/auth"

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
		app.authHandle,
		app.userHandle,
		app.roleHandle,
		app.coachHandle,
		app.coacheeHandle,
	}
}

func initApp(cfg *configs.Config) (*AppApi, error) {
	dbConn := dbInitConnection(cfg)

	userRepo := repouser.NewUserRepo(dbConn)
	roleRepo := reporole.NewRoleRepo(dbConn)

	roleService := servicerole.NewRoleService(roleRepo)
	userService := serviceuser.NewUserService(userRepo)
	authService := serviceauth.NewAuthService(userRepo)
	coachService := servicecoach.NewCoachService(userRepo)
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
