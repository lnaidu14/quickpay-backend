package main

import (
	"net/http"

	"quickpay/main/helpers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Health check to see if server is live
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Server is live")

	})

	// Returning a user name
	app.Get("/api/user/:id", func(c *fiber.Ctx) error {
		helpers.GenQrCode(helpers.ExampleUser)
		return c.Status(http.StatusOK).JSON(helpers.ExampleUser)
	})

	app.Listen(":3000")
}
