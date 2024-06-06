package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Person struct {
	Name string `json:"name"`
}

func main() {
	app := fiber.New()

	response := Person{"Lalit"}

	// Health check to see if server is live
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Server is live")

	})

	// Returning a user name
	app.Get("/api/user", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(response)
	})

	app.Listen(":3000")
}
