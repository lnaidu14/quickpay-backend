package types

import "time"

type UserTransactions struct {
	Tx_Id       string    `json:"tx_id"`
	Amt         int       `json:"amt"`
	Tx_Datetime time.Time `json:"tx_datetime"`
}

type UserTransactionBody struct {
	Amt uint `json:"amt"`
}
