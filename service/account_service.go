package service

import (
	"fmt"
	"internal-transfers/common"
	"internal-transfers/model"
	"internal-transfers/persistence"
	"time"

	"context"

	"github.com/shopspring/decimal"
)

type AccountService struct {
	Repo        *persistence.AccountRepository
	AuditLogger *common.AuditLogger
}

func NewAccountService(accountRepo *persistence.AccountRepository, auditLogger *common.AuditLogger) *AccountService {
	return &AccountService{
		Repo:        accountRepo,
		AuditLogger: auditLogger,
	}
}

// Creates a new account with retry mechanism for database errors
func (accountService *AccountService) CreateAccount(account model.Account) error {
	var err error
	for i := 0; i < 2; i++ {
		err = accountService.createAccountWithRetry(account)
		if err == nil {
			return nil
		}

		// Retry if it's a database error or timeout
		if err.Error() == "context deadline exceeded" || err.Error() == "pq: deadlock detected" {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	return err
}

// createAccount with Retry actually handles the creation with context and timeout context 30 seconds
func (accountService *AccountService) createAccountWithRetry(account model.Account) error {
	// Set a timeout context (30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	existingAccount, err := accountService.Repo.GetAccountByIDWithContext(ctx, account.AccountID)
	if err != nil {
		return fmt.Errorf("error checking if account exists: %v", err)
	}

	if existingAccount != nil {
		return fmt.Errorf("account already exists")
	}

	err = accountService.Repo.CreateAccountWithContext(ctx, account)
	if err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}

	return nil
}

// Retrieves an account by its ID with retry mechanism
func (accountService *AccountService) GetAccountByID(accountID int) (*model.Account, error) {
	var err error
	for i := 0; i < 2; i++ {
		account, err := accountService.getAccountByIDWithRetry(accountID)
		if err == nil {
			return account, nil
		}

		if err.Error() == "context deadline exceeded" || err.Error() == "pq: deadlock detected" {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	return nil, err
}

// Handles the retrieval with context and timeout
func (accountService *AccountService) getAccountByIDWithRetry(accountID int) (*model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	account, err := accountService.Repo.GetAccountByIDWithContext(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account by ID: %v", err)
	}

	return account, nil
}

// Updates the balance of an existing account with retry mechanism
func (accountService *AccountService) UpdateAccountBalance(accountID int, newBalance decimal.Decimal) error {
	var err error
	for i := 0; i < 2; i++ {
		err = accountService.updateAccountBalanceWithRetry(accountID, newBalance)
		if err == nil {
			return nil
		}

		if err.Error() == "context deadline exceeded" || err.Error() == "pq: deadlock detected" {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	return err
}

// Handles the balance update with context and timeout
func (accountService *AccountService) updateAccountBalanceWithRetry(accountID int, newBalance decimal.Decimal) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := accountService.Repo.UpdateAccountBalanceWithContext(ctx, accountID, newBalance)
	if err != nil {
		return fmt.Errorf("error updating account balance: %v", err)
	}

	return nil
}
