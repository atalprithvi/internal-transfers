package model

import "github.com/shopspring/decimal"

type Account struct {
	AccountID int             `json:"account_id" db:"account_id"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
}

func NewAccount(accountID int, initialBalance decimal.Decimal) *Account {
	return &Account{
		AccountID: accountID,
		Balance:   initialBalance,
	}
}

type CreateAccountInput struct {
	AccountID      int             `json:"account_id"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
}

type UpdateAccountInput struct {
	AccountID int             `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
}
