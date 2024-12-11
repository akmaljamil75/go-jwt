package router

import (
	"go-jwt/modules/auth"
	"go-jwt/modules/role"
	"go-jwt/modules/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRouterPrivate(db *gorm.DB, c *fiber.App) *fiber.App {

	api := c.Group("/api/v1")

	// Role ROUTER API
	roleRoute := api.Group("/role")
	roleRepository := role.NewRepository(db)
	roleService := role.NewService(roleRepository)
	roleHandler := role.NewHandler(roleService)
	roleRoute.Post("/create", roleHandler.Create)
	roleRoute.Post("/search", roleHandler.FindRoles)
	roleRoute.Get("/:name", roleHandler.FindOneRoleByName)
	roleRoute.Get("/:id", roleHandler.FindOneRoleByID)
	roleRoute.Patch("/:id", roleHandler.Update)
	roleRoute.Delete("/:id", roleHandler.SoftDelete)
	roleRoute.Put("/:id", roleHandler.RestoreSoftDelete)

	// USER ROUTER API
	userRoute := api.Group("/user")
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, roleRepository)
	userHandler := user.NewHandler(userService)
	userRoute.Post("/create", userHandler.Create)

	return c
}

func InitRouterPublic(db *gorm.DB, c *fiber.App) *fiber.App {

	api := c.Group("/api/auth")

	roleRepository := role.NewRepository(db)
	userRepository := user.NewRepository(db)

	authService := auth.NewService(userRepository, roleRepository)
	authHandler := auth.NewHandler(authService)

	api.Post("/login", authHandler.Login)

	return c
}
