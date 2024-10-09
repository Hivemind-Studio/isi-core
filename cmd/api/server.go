package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	handleuser "github.com/Hivemind-Studio/isi-core/internal/handler/http/user"
	repouser "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	serviceuser "github.com/Hivemind-Studio/isi-core/internal/service/user"

	handlerole "github.com/Hivemind-Studio/isi-core/internal/handler/http/role"
	reporole "github.com/Hivemind-Studio/isi-core/internal/repository/role"
	servicerole "github.com/Hivemind-Studio/isi-core/internal/service/role"

	"github.com/gofiber/fiber/v2"
)

type AppApi struct {
	userHandle *handleuser.Handler
	roleHandle *handlerole.Handler
}

type Router interface {
	RegisterRoutes(app *fiber.App)
}

func routerList(app *AppApi) []Router {
	return []Router{
		app.userHandle,
		app.roleHandle,
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

	return &AppApi{userHandle: userHandler, roleHandle: roleHandler}, nil
}
