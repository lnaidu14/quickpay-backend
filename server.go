package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"quickpay/main/helpers"
	"quickpay/main/types"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()

	// Health check to see if server is live
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Server is live")
	})

	// Fetch and check if user exists
	app.Get("/api/users/:id", func(c *fiber.Ctx) error {
		url := fmt.Sprintf("https://dev-quickpay.us.auth0.com/api/v2/users/%s", c.Params("id"))

		request := fiber.Get(url)

		// to set headers
		request.Set("'Accept", "application/json")
		request.Set("authorization", os.Getenv("AUTH0_MGMT_API_TOKEN"))

		statusCode, data, errs := request.Bytes()

		if len(errs) > 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"errs": errs,
			})
		}

		var value fiber.Map
		err := json.Unmarshal(data, &value)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": err,
			})
		}

		return c.Status(statusCode).JSON(value)
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
