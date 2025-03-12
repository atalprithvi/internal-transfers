package model

import "github.com/shopspring/decimal"

type Transaction struct {
	SourceAccountID      int             `json:"source_account_id" db:"source_account_id"`
	DestinationAccountID int             `json:"destination_account_id" db:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount" db:"amount"`
}

type TransactionRequest struct {
	SourceAccountID      int             `json:"source_account_id"`
	DestinationAccountID int             `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
}

func NewTransaction(sourceAccountID, destinationAccountID int, amount decimal.Decimal) *Transaction {
	return &Transaction{
		SourceAccountID:      sourceAccountID,
		DestinationAccountID: destinationAccountID,
		Amount:               amount,
	}
}
