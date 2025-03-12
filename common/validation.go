package common

import (
	"errors"
	"internal-transfers/model"

	"github.com/shopspring/decimal"
)

// Validate if the account balance is sufficient for the transaction
func ValidateBalance(account model.Account, transactionAmount decimal.Decimal) error {
	if account.Balance.LessThan(transactionAmount) {
		return errors.New("insufficient balance")
	}
	return nil
}

// Validate if the account already exists
func ValidateAccountExistence(accountID int, accounts []model.Account) error {
	for _, acc := range accounts {
		if acc.AccountID == accountID {
			return nil
		}
	}
	return errors.New("account does not exist")
}
