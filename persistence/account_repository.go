package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"internal-transfers/model"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// Defines methods to interact with the accounts table in the database
type AccountRepository struct {
	DB *sqlx.DB
}

// NewAccountRepository creates a new AccountRepository
func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

// Retrieves an account by its ID using context with timeout
func (repo *AccountRepository) GetAccountByIDWithContext(ctx context.Context, accountID int) (*model.Account, error) {
	var account model.Account
	query := `SELECT account_id, balance FROM accounts WHERE account_id = $1`
	err := repo.DB.GetContext(ctx, &account, query, accountID)
	if err != nil {
		if err == sql.ErrNoRows { 
			return nil, nil 
		}
		return nil, fmt.Errorf("error getting account by ID: %v", err)
	}
	return &account, nil
}

// Creates a new account using context with timeout
func (repo *AccountRepository) CreateAccountWithContext(ctx context.Context, account model.Account) error {
	query := `INSERT INTO accounts (account_id, balance) VALUES ($1, $2)`
	_, err := repo.DB.ExecContext(ctx, query, account.AccountID, account.Balance.String())
	if err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}
	return nil
}

// Updates the balance of an existing account using context with timeout
func (repo *AccountRepository) UpdateAccountBalanceWithContext(ctx context.Context, accountID int, newBalance decimal.Decimal) error {
	query := `UPDATE accounts SET balance = $1 WHERE account_id = $2`
	_, err := repo.DB.ExecContext(ctx, query, newBalance.String(), accountID)
	if err != nil {
		return fmt.Errorf("error updating account balance: %v", err)
	}
	return nil
}
