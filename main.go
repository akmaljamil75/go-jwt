package main

import (
	"fmt"
	"go-jwt/common/database"
	"go-jwt/common/middleware"
	"go-jwt/common/response"
	"go-jwt/common/router"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func catch() {
	if r := recover(); r != nil {
		fmt.Println("Error occured", r)
	} else {
		fmt.Println("Application running perfectly")
	}
}

func main() {
	defer catch()
	db := database.InitDB()
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(middleware.LoggerMiddleware)
	app.Use(middleware.HandlingErrorMiddleware)
	app = router.InitRouterPublic(db, app)
	app.Use(middleware.JwtAuthorization)
	app = router.InitRouterPrivate(db, app)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(response.BuildFailedResponseMessage("service not found", 404, nil))
	})
	app.Listen(":8081")
}
