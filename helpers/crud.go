package helpers

import (
	"context"
	"fmt"
	"log"
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
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
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

func FetchUserBalance(conn *pgx.Conn, userId string) (int, error) {
	log.Println("Entering FetchUserBalance()")
	var userBalance int
	err := conn.QueryRow(context.Background(), "SELECT balance FROM users WHERE id=$1", userId).Scan(&userBalance)
	if err != nil {
		log.Printf("Error occured when querying user balance: %v\n", err)
		return -1, err
	}
	log.Println("Exiting FetchUserBalance()")
	return userBalance, nil
}

func FetchUserTransactions(conn *pgx.Conn, userId string) (*[]types.UserTransactions, error) {
	log.Println("Entering FetchUserTransactions()")
	rows, err := conn.Query(context.Background(), "SELECT tx_id, amt, tx_datetime FROM transactions INNER JOIN users ON users.id = transactions.user_id WHERE id=$1 ORDER BY tx_datetime DESC LIMIT 10", userId)
	if err != nil {
		log.Printf("Error occured when querying user transactions: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	userTransactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.UserTransactions])
	if err != nil {
		log.Printf("Error occured when iterating rows of user transactions: %v\n", err)
		return &userTransactions, err
	}
	log.Println("Exiting FetchUserTransactions()")
	return &userTransactions, nil
}

func CreateUserTransaction(conn *pgx.Conn, userId string, payload types.UserTransactionBody) (types.Response, error) {
	log.Println("Entering CreateUserTransaction()")

	id := uuid.New()

	_, err := conn.Exec(context.Background(), "INSERT INTO transactions (tx_id, user_id, amt, tx_datetime) values ($1, $2, $3, $4);", id.String(), userId, payload.Amt, time.Now())
	if err != nil {
		log.Printf("Error occured when creating user transactions: %v\n", err)
		return types.Response{
			Message: "Error occured when creating user transactions",
		}, err
	}

	log.Println("Exiting CreateUserTransaction()")
	return types.Response{
		Message: "Successfully created user transaction",
	}, nil
}

func UpdateUserBalance(conn *pgx.Conn, userId string, amt types.UserTransactionBody) error {
	log.Println("Entering UpdateUserBalance()")

	currUserBalance, err := FetchUserBalance(conn, userId)
	if err != nil {
		log.Printf("Error occured when fetching user balance: %v\n", err)
		return nil
	}

	updatedUserBalance := currUserBalance - amt.Amt

	_, err = conn.Exec(context.Background(), "UPDATE users SET balance=$1 WHERE id=$2;", updatedUserBalance, userId)
	if err != nil {
		log.Printf("Error occured when updating user balance: %v\n", err)
		return err
	}

	log.Println("Exiting UpdateUserBalance()")
	return nil
}
