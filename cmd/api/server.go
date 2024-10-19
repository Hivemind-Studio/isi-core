package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	handleuser "github.com/Hivemind-Studio/isi-core/internal/handler/http/user"
	repouser "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	serviceuser "github.com/Hivemind-Studio/isi-core/internal/service/user"

	handlerole "github.com/Hivemind-Studio/isi-core/internal/handler/http/role"
	reporole "github.com/Hivemind-Studio/isi-core/internal/repository/role"
	servicerole "github.com/Hivemind-Studio/isi-core/internal/service/role"

	handleAuth "github.com/Hivemind-Studio/isi-core/internal/handler/http/auth"
	repoAuth "github.com/Hivemind-Studio/isi-core/internal/repository/auth"
	serviceAuth "github.com/Hivemind-Studio/isi-core/internal/service/auth"

	"github.com/gofiber/fiber/v2"
)

type AppApi struct {
	userHandle *handleuser.UserHandler
	roleHandle *handlerole.RoleHandler
	authHandle *handleAuth.AuthHandler
}

type Router interface {
	RegisterRoutes(app *fiber.App)
}

func routerList(app *AppApi) []Router {
	return []Router{
		app.userHandle,
		app.roleHandle,
		app.authHandle,
	}
}

func initApp(cfg *configs.Config) (*AppApi, error) {
	dbConn := dbInitConnection(cfg)

	userRepo := repouser.NewUserRepo(dbConn)
	userService := serviceuser.NewUserService(userRepo)
	userHandler := handleuser.NewUserHandler(userService)

	roleRepo := reporole.NewRoleRepo(dbConn)
	roleService := servicerole.NewRoleService(roleRepo)
	roleHandler := handlerole.NewRoleHandler(roleService)

	authRepo := repoAuth.NewAuthRepo(dbConn)
	authService := serviceAuth.NewAuthService(authRepo)
	authHandler := handleAuth.NewAuthHandler(authService)

	return &AppApi{userHandle: userHandler,
		roleHandle: roleHandler,
		authHandle: authHandler}, nil
}
