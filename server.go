package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"quickpay/main/helpers"
	"quickpay/main/types"
)

func main() {
	app := fiber.New()

	// Health check to see if server is live
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Server is live")

	})

	// Returning a user name
	app.Post("/api/user/:id", func(c *fiber.Ctx) error {

		var payload types.User

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).SendString("Error occured when generating QR code")
		}

		imageBase64String := helpers.GenQrCode(payload)
		return c.Status(http.StatusOK).SendString(imageBase64String)
	})

	app.Listen(":3000")
}
