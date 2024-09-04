package types

import "time"

type UserTransactions struct {
	Tx_Id       string    `json:"tx_id"`
	SenderId    string    `json:"sender_id"`
	UserId      string    `json:"user_id"`
	Amt         int       `json:"amt"`
	Tx_Datetime time.Time `json:"tx_datetime"`
}

type UserTransactionBody struct {
	Amt         int    `json:"amt"`
	RecipientId string `json:"recipientId"`
	SenderId    string `json:"senderId"`
}
