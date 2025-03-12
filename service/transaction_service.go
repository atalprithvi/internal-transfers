package service

import (
	"context"
	"fmt"
	"internal-transfers/common"
	"internal-transfers/model"
	"internal-transfers/persistence"
	"time"

	"github.com/shopspring/decimal"
)

// Responsible for handling the transaction related business logic
type TransactionService struct {
	AccountRepo     *persistence.AccountRepository
	TransactionRepo *persistence.TransactionRepository
	AuditLogger     *common.AuditLogger
}

func NewTransactionService(accountRepo *persistence.AccountRepository, transactionRepo *persistence.TransactionRepository, auditLogger *common.AuditLogger) *TransactionService {
	return &TransactionService{
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
		AuditLogger:     auditLogger,
	}
}

// Handles the logic for performing a transaction with retry and timeout
func (transactionService *TransactionService) PerformTransaction(transaction model.Transaction) error {
	// Retry mechanism
	var err error
	for i := 0; i < 2; i++ {
		err = transactionService.performTransactionWithRetry(transaction)
		if err == nil {
			return nil
		}

		if err.Error() == "context deadline exceeded" || err.Error() == "pq: deadlock detected" {
			time.Sleep(2 * time.Second) // wait for 2 seconds before retrying
			continue
		}
		break
	}
	return err
}

// performTransactionWithRetry actually handles the transaction with context and timeout
func (transactionService *TransactionService) performTransactionWithRetry(transaction model.Transaction) error {
	// Set a timeout context (30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if transactionService.AuditLogger != nil {
		transactionService.AuditLogger.LogAction("Transaction Initiated", fmt.Sprintf("Source Account ID: %d, Destination Account ID: %d, Amount: %s",
			transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount.String()))
	}

	if transaction.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return fmt.Errorf("transaction amount must be greater than zero")
	}

	if transactionService.AuditLogger != nil {
		transactionService.AuditLogger.LogAction("Amount Validation", fmt.Sprintf("Transaction Amount: %s (Valid: %t)", transaction.Amount.String(), transaction.Amount.GreaterThan(decimal.NewFromInt(0))))
	}

	sourceAccount, err := transactionService.AccountRepo.GetAccountByIDWithContext(ctx, transaction.SourceAccountID)
	if err != nil {
		return fmt.Errorf("source account validation failed: %v", err)
	}

	if sourceAccount.Balance.LessThan(transaction.Amount) {
		return fmt.Errorf("insufficient balance in source account")
	}

	destinationAccount, err := transactionService.AccountRepo.GetAccountByIDWithContext(ctx, transaction.DestinationAccountID)
	if err != nil {
		return fmt.Errorf("destination account validation failed: %v", err)
	}

	if transactionService.AuditLogger != nil {
		transactionService.AuditLogger.LogAction("Destination Account Found", fmt.Sprintf("Destination Account ID: %d, Balance: %s", destinationAccount.AccountID, destinationAccount.Balance.String()))
	}

	if err := transactionService.updateAccountBalanceWithContext(ctx, sourceAccount, transaction.Amount.Neg()); err != nil {
		return err
	}

	if err := transactionService.updateAccountBalanceWithContext(ctx, destinationAccount, transaction.Amount); err != nil {
		transactionService.updateAccountBalanceWithContext(ctx, sourceAccount, transaction.Amount)
		return err
	}

	// Save the transaction record
	if err := transactionService.TransactionRepo.SaveTransactionWithContext(ctx, transaction); err != nil {
		transactionService.updateAccountBalanceWithContext(ctx, sourceAccount, transaction.Amount)
		transactionService.updateAccountBalanceWithContext(ctx, destinationAccount, transaction.Amount.Neg())
		return fmt.Errorf("failed to save transaction: %v", err)
	}

	if transactionService.AuditLogger != nil {
		transactionService.AuditLogger.LogAction("Transaction Completed", fmt.Sprintf("Transaction from Account %d to Account %d for Amount: %s",
			transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount.String()))
	}

	return nil
}

// updateAccountBalanceWithContext updates the account balance with retry and context
func (transactionService *TransactionService) updateAccountBalanceWithContext(ctx context.Context, account *model.Account, amount decimal.Decimal) error {
	account.Balance = account.Balance.Add(amount)
	if err := transactionService.AccountRepo.UpdateAccountBalanceWithContext(ctx, account.AccountID, account.Balance); err != nil {
		return fmt.Errorf("failed to update account balance for Account ID %d: %v", account.AccountID, err)
	}
	return nil
}
