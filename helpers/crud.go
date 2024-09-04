package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"quickpay/main/types"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
)

func NewConnection() (*pgx.Conn, error) {
	log.Println("Entering NewConnection()...")
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewConnect(): Unable to connect to database: %v\n", err)
		return nil, &types.ErrorResponse{Err: errors.New("unable to connect to database")}
	}
	log.Println("Connection established...")
	log.Println("Exiting NewConnection()...")
	return conn, nil
}

// func FetchUsers(conn *pgx.Conn) (*[]types.UserBalance, error) {
// 	log.Println("Entering FetchUsers()")
// 	rows, err := conn.Query(context.Background(), "SELECT * FROM users")
// 	if err != nil {
// 		log.Printf("Error occured when querying users: %v\n", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.UserBalance])
// 	if err != nil {
// 		log.Printf("Error occured when iterating rows of users: %v\n", err)
// 		return &users, err
// 	}

// 	log.Println("Exiting FetchUsers()")
// 	return &users, nil
// }

func FetchUserFromDb(conn *pgx.Conn, userId string) (*types.DbUser, error) {
	log.Println("Entering FetchUserFromDb()")
	var user types.DbUser
	decodedUserId, err := url.QueryUnescape(userId)
	if err != nil {
		panic(err)
	}
	fmt.Println("decodedUserId: ", decodedUserId)
	err = conn.QueryRow(context.Background(), "SELECT * FROM users WHERE id=$1", decodedUserId).Scan(&user.Id, &user.Balance, &user.Username)
	if err != nil {
		log.Printf("FetchUserFromDb(): Error occured when querying for user: %v\n", err)
		return &types.DbUser{}, &types.ErrorResponse{Err: errors.New("error when querying for user")}
	}
	log.Println("Exiting FetchUserFromDb()")
	return &user, nil
}

func CreateUserInDb(conn *pgx.Conn, payload types.CreateDbUser) (*types.GenericResponse, error) {
	log.Println("Entering CreateUserInDb()")
	_, err := conn.Exec(context.Background(), "INSERT INTO users (id, balance, username) values ($1, 0, $2);", payload.Id, payload.Username)
	if err != nil {
		log.Printf("CreateUserInDb(): Error occured when creating user: %v\n", err)
		return nil, &types.ErrorResponse{Err: errors.New("error when creating user")}
	}

	log.Println("Exiting CreateUserInDb()")
	return &types.GenericResponse{
		Message: "Successfully created user",
	}, nil
}

func FetchUserBalance(conn *pgx.Conn, userId string) (int, error) {
	log.Println("Entering FetchUserBalance()")
	var userBalance int
	decodedUserId, err := url.QueryUnescape(userId)
	if err != nil {
		panic(err)
	}
	err = conn.QueryRow(context.Background(), "SELECT balance FROM users WHERE id=$1", decodedUserId).Scan(&userBalance)
	if err != nil {
		log.Printf("FetchUserBalance(): Error occured when querying user balance: %v\n", err)
		return -1, &types.ErrorResponse{Err: errors.New("error occured when querying user balance")}
	}
	log.Println("Exiting FetchUserBalance()")
	return userBalance, nil
}

func FetchUserTransactions(conn *pgx.Conn, userId string) (*[]types.UserTransactions, error) {
	log.Println("Entering FetchUserTransactions()")
	decodedUserId, err := url.QueryUnescape(userId)
	if err != nil {
		panic(err)
	}
	rows, err := conn.Query(context.Background(), "SELECT tx_id, amt, tx_datetime, user_id, sender_id FROM transactions INNER JOIN users ON users.id = transactions.user_id WHERE id=$1 ORDER BY tx_datetime DESC LIMIT 10", decodedUserId)
	if err != nil {
		log.Printf("FetchUserBalance(): Error occured when querying user transactions: %v\n", err)
		return nil, &types.ErrorResponse{Err: errors.New("error occured when querying user transactions")}
	}
	defer rows.Close()

	userTransactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.UserTransactions])
	if err != nil {
		log.Printf("FetchUserTransactions(): Error occured when iterating rows of user transactions: %v\n", err)
		return &userTransactions, &types.ErrorResponse{Err: errors.New("error occured when iterating rows of user transactions")}
	}
	log.Println("Exiting FetchUserTransactions()")
	return &userTransactions, nil
}

func CreateUserTransaction(conn *pgx.Conn, payload types.UserTransactionBody) (*types.GenericResponse, error) {
	log.Println("Entering CreateUserTransaction()")

	fmt.Println("Payload: ", payload)

	decodedSenderUserId, err := url.QueryUnescape(payload.SenderId)
	if err != nil {
		panic(err)
	}

	decodedRecipientUserId, err := url.QueryUnescape(payload.RecipientId)
	if err != nil {
		panic(err)
	}

	senderUserName, err := FetchUserFromDb(conn, decodedSenderUserId)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO transactions (tx_id, user_id, amt, tx_datetime, sender_id) values ($1, $2, $3, $4, $5);", uuid.New().String(), decodedSenderUserId, payload.Amt, time.Now(), "N/A")

	if err != nil {
		log.Printf("CreateUserTransaction(): Error occured when creating sender user transactions: %v\n", err)
		return nil, &types.ErrorResponse{Err: errors.New("error occured when creating user transactions")}
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO transactions (tx_id, user_id, amt, tx_datetime, sender_id) values ($1, $2, $3, $4, $5);", uuid.New().String(), decodedRecipientUserId, payload.Amt, time.Now(), senderUserName.Username)

	if err != nil {
		log.Printf("CreateUserTransaction(): Error occured when creating user transactions: %v\n", err)
		return nil, &types.ErrorResponse{Err: errors.New("error occured when creating recipient user transactions")}
	}

	log.Println("Exiting CreateUserTransaction()")
	return &types.GenericResponse{
		Message: "Successfully created user transaction",
	}, nil
}

func UpdateUserBalance(conn *pgx.Conn, transactionBody types.UserTransactionBody) error {
	log.Println("Entering UpdateUserBalance()")

	decodedSenderUserId, err := url.QueryUnescape(transactionBody.SenderId)
	if err != nil {
		panic(err)
	}

	senderUserBalance, err := FetchUserBalance(conn, decodedSenderUserId)
	if err != nil {
		log.Printf("UpdateUserBalance(): Error occured when fetching sender user balance: %v\n", err)
		return &types.ErrorResponse{Err: errors.New("error occured when fetching sender user balance")}
	}

	if senderUserBalance < 1 || senderUserBalance < transactionBody.Amt {
		log.Printf("UpdateUserBalance(): Insufficient funds")
		return &types.ErrorResponse{Err: errors.New("insufficient funds")}
	}

	updatedSenderBalance := senderUserBalance - transactionBody.Amt

	decodedRecipientUserId, err := url.QueryUnescape(transactionBody.RecipientId)
	if err != nil {
		panic(err)
	}

	recipientUserBalance, err := FetchUserBalance(conn, decodedRecipientUserId)
	if err != nil {
		log.Printf("UpdateUserBalance(): Error occured when fetching recipient user balance: %v\n", err)
		return &types.ErrorResponse{Err: errors.New("error occured when fetching recipient user balance")}
	}

	updatedRecipientBalance := recipientUserBalance + transactionBody.Amt

	_, err = conn.Exec(context.Background(), "UPDATE users SET balance=$1 WHERE id=$2;", updatedSenderBalance, decodedSenderUserId)
	if err != nil {
		log.Printf("UpdateUserBalance(): Error occured when updating sender balances: %v\n", err)
		return &types.ErrorResponse{Err: errors.New("error occured when updating balance")}
	}
	// TODO: Figure out how to update in one transaction
	_, err = conn.Exec(context.Background(), "UPDATE users SET balance=$1 WHERE id=$2;", updatedRecipientBalance, decodedRecipientUserId)
	if err != nil {
		log.Printf("UpdateUserBalance(): Error occured when updating recipient balances: %v\n", err)
		return &types.ErrorResponse{Err: errors.New("error occured when updating balance")}
	}

	log.Println("Exiting UpdateUserBalance()")
	return nil
}
