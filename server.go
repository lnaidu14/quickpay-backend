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

	// Fetch user balance
	app.Get("/api/users/:id/balance", func(c *fiber.Ctx) error {
		conn, err := helpers.NewConnection()

		if err != nil {
			log.Printf("Error occured when connecting to database: %v\n", err)
			return nil
		}

		userId := c.Params("id")

		userBalance, err := helpers.FetchUserBalance(conn, userId)
		if err != nil {
			log.Printf("Error occured when fetching user balance: %v\n", err)
			return nil
		}
		return c.Status(http.StatusOK).JSON(userBalance)
	})

	// Fetch user transactions
	app.Get("/api/users/:id/transactions", func(c *fiber.Ctx) error {
		conn, err := helpers.NewConnection()

		if err != nil {
			log.Printf("Error occured when connecting to database: %v\n", err)
			return nil
		}

		userId := c.Params("id")

		userTransactions, err := helpers.FetchUserTransactions(conn, userId)
		if err != nil {
			log.Printf("Error occured when fetching user transactions: %v\n", err)
			return nil
		}
		return c.Status(http.StatusOK).JSON(userTransactions)
	})

	// Create user transaction
	app.Post("/api/users/:id/transactions", func(c *fiber.Ctx) error {
		conn, err := helpers.NewConnection()

		if err != nil {
			log.Printf("Error occured when connecting to database: %v\n", err)
			return nil
		}

		userId := c.Params("id")
		var payload types.UserTransactionBody
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).SendString("Error occured when creating transaction")
		}

		// Create transaction record
		message, err := helpers.CreateUserTransaction(conn, userId, payload)
		if err != nil {
			log.Printf("Error occured when fetching user transactions: %v\n", err)
			return nil
		}
		// Update balance of user
		return c.Status(http.StatusAccepted).JSON(message)
	})

	// app.Get("/api/users", func(c *fiber.Ctx) error {
	// 	// TODO: Fetch users from Auth0 directly after adding payments, user payments etc.
	// 	// urlExample := "postgres://username:password@localhost:5432/database_name"
	// 	conn, err := helpers.NewConnection()
	// 	if err != nil {
	// 		log.Printf("Error occured when connecting to database: %v\n", err)
	// 		return nil
	// 	}

	// 	users, err := helpers.FetchUsers(conn)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 		return c.Status(http.StatusBadRequest).SendString("Query failed")
	// 	}
	// 	return c.Status(http.StatusOK).JSON(users)
	// })

	// Fetch and check if user exists in Auth0
	app.Get("/api/users/:id", func(c *fiber.Ctx) error {
		// cipher key
		key := os.Getenv("CIPHER_KEY")

		encryptedUserData := c.Params("id")

		// decrypt
		decryptedUserData, _ := helpers.DecryptAES([]byte(key), encryptedUserData)
		fmt.Println("decrypted: ", decryptedUserData)

		var msg types.User
		err := json.Unmarshal([]byte(decryptedUserData), &msg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": err,
			})
		}

		url := fmt.Sprintf("https://dev-quickpay.us.auth0.com/api/v2/users/%s", msg.Id)

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
		unmarshalError := json.Unmarshal(data, &value)
		if unmarshalError != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": unmarshalError,
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
