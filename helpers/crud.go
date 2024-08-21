package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"quickpay/main/types"

	"github.com/jackc/pgx/v5"
)

func NewConnection() (*pgx.Conn, error) {
	log.Println("Entering NewConnection()...")
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())
	log.Println("Exiting NewConnection()...")
	return conn, nil
}

func FetchUsers(conn *pgx.Conn) (*[]types.User, error) {
	log.Println("Entering FetchUsers()")
	rows, err := conn.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		log.Printf("Error occured when fetching users from database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.User])
	if err != nil {
		log.Printf("Error occured when iterating rows: %v\n", err)
		return &users, err
	}

	log.Println("Exiting FetchUsers()")
	return &users, nil
}
