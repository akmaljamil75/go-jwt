package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(c *fiber.Ctx) error {

	log.Printf("Request: %s %s", c.Method(), c.Path())

	err := c.Next()
	if err != nil {
		log.Printf("Error: %v", err)
	}

	log.Printf("Response: %d %s", c.Response().StatusCode(), string(c.Response().Body()))
	return nil
}
