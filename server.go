package main

import (
	"context"
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
			return c.SendStatus(http.StatusInternalServerError)
		}

		defer conn.Close(context.Background())

		userId := c.Params("id")

		userBalance, err := helpers.FetchUserBalance(conn, userId)
		if err != nil {
			log.Printf("Error occured when fetching user balance: %v\n", err)
			return err
		}
		return c.Status(http.StatusOK).JSON(userBalance)
	})

	// Fetch user from database
	app.Get("/api/users/:id", func(c *fiber.Ctx) error {
		conn, err := helpers.NewConnection()

		if err != nil {
			log.Printf("Error occured when connecting to database: %v\n", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		defer conn.Close(context.Background())

		userId := c.Params("id")

		user, err := helpers.FetchUserFromDb(conn, userId)
		if err != nil {
			log.Printf("Error occured when fetching user from database: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}
		return c.Status(http.StatusOK).JSON(user)
	})

	// Fetch user transactions
	app.Get("/api/users/:id/transactions", func(c *fiber.Ctx) error {
		conn, err := helpers.NewConnection()

		if err != nil {
			log.Printf("Error occured when connecting to database: %v\n", err)
			return nil
		}

		defer conn.Close(context.Background())

		userId := c.Params("id")

		userTransactions, err := helpers.FetchUserTransactions(conn, userId)
		if err != nil {
			log.Printf("Error occured when fetching user transactions: %v\n", err)
			return nil
		}
		return c.Status(http.StatusOK).JSON(userTransactions)
	})

	// Create user transaction
	app.Post("/api/users/transactions", func(c *fiber.Ctx) error {
		log.Printf("Creating user transaction for users")
		conn, err := helpers.NewConnection()

		if err != nil {
			log.Printf("Error occured when connecting to database: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}

		defer conn.Close(context.Background())

		var payload types.UserTransactionBody
		if err := c.BodyParser(&payload); err != nil {
			log.Printf("Error occured when parsing payload: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}

		if payload.Amt == 0 {
			log.Printf("Amount transferred cannot be zero")
			return c.Status(http.StatusBadRequest).JSON(types.GenericResponse{Message: "amount transferred cannot be zero"})
		}

		// Update balance of user
		err = helpers.UpdateUserBalance(conn, payload)
		if err != nil {
			log.Printf("Error occured when updating user balance: %v\n", err)
			return c.Status(http.StatusBadRequest).JSON(types.GenericResponse{Message: err.Error()})
		}

		// Create transaction record
		transactionMessage, err := helpers.CreateUserTransaction(conn, payload)
		if err != nil {
			log.Printf("Error occured creating user transactions: %v\n", err)
			return c.Status(http.StatusBadRequest).JSON(types.GenericResponse{
				Message: err.Error(),
			})
		}

		return c.Status(http.StatusCreated).JSON(transactionMessage)
	})

	// app.Get("/api/users", func(c *fiber.Ctx) error {
	// 	// TODO: Fetch users from Auth0 directly after adding payments, user payments etc.
	// 	// urlExample := "postgres://username:password@localhost:5432/database_name"
	// 	conn, err := helpers.NewConnection()
	// 	if err != nil {
	// 		log.Printf("Error occured when connecting to database: %v\n", err)
	// 		return nil
	// 	}

	//  defer conn.Close(context.Background())

	// 	users, err := helpers.FetchUsers(conn)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 		return c.Status(http.StatusBadRequest).SendString("Query failed")
	// 	}
	// 	return c.Status(http.StatusOK).JSON(users)
	// })

	// Fetch and check if user exists in Auth0
	app.Get("/api/auth/users/:id", func(c *fiber.Ctx) error {
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

	// Generate a user QR code
	app.Post("/api/users/qr/:id", func(c *fiber.Ctx) error {
		fmt.Println("Generating QR code...")
		var payload types.User

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).SendString("Error occured when generating QR code")
		}

		imageBase64String := helpers.GenQrCode(payload)
		return c.Status(http.StatusOK).SendString(imageBase64String)
	})

	// Create user
	app.Post("/api/users", func(c *fiber.Ctx) error {
		conn, err := helpers.NewConnection()

		if err != nil {
			log.Printf("Error occured when connecting to database: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}

		defer conn.Close(context.Background())

		var payload types.CreateDbUser

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).SendString("Error occured when generating QR code")
		}

		message, err := helpers.CreateUserInDb(conn, payload)

		if err != nil {
			log.Printf("Error occured when fetching user transactions: %v\n", err)
			return c.Status(http.StatusBadRequest).JSON(err.Error())

		}

		return c.Status(http.StatusOK).JSON(message)
	})

	app.Post("/api/authorize/:id", func(c *fiber.Ctx) error {

		userId := c.Params("id")

		response, err := helpers.FetchManagementApiToken()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": err,
			})
		}
		accessToken := response.AccessToken

		profile, err := helpers.FetchUserProfile(accessToken, userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": err,
			})
		}
		fmt.Println("profile: ", profile)
		return c.Status(http.StatusOK).JSON(profile)

	})

	app.Listen(":3000")
}
