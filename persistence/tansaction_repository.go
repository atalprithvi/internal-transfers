package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"internal-transfers/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Responsible for interacting with the database for transaction related operations
type TransactionRepository struct {
	DB *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

// Retrieves a transaction by its ID
func (transactionRepository *TransactionRepository) GetTransactionByID(transactionID string) (*model.Transaction, error) {
	query := `SELECT transaction_id, source_account_id, destination_account_id, amount 
			  FROM transactions 
			  WHERE transaction_id = $1`

	var transaction model.Transaction

	err := transactionRepository.DB.Get(&transaction, query, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with ID %s not found", transactionID)
		}
		return nil, fmt.Errorf("failed to fetch transaction: %v", err)
	}

	return &transaction, nil
}

// SaveTransaction saves a transaction to the database (without context support)
func (transactionRepository *TransactionRepository) SaveTransaction(transaction model.Transaction) error {
	query := `INSERT INTO transactions (source_account_id, destination_account_id, amount)
	VALUES ($1, $2, $3)`

	_, err := transactionRepository.DB.Exec(query, transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount.String())
	if err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}

	return nil
}

// Saves a transaction to the database with context support for timeout and cancellation
func (transactionRepository *TransactionRepository) SaveTransactionWithContext(ctx context.Context, transaction model.Transaction) error {
	tx, err := transactionRepository.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO transactions (source_account_id, destination_account_id, amount)
	VALUES ($1, $2, $3)`

	_, err = tx.ExecContext(ctx, query, transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount.String())
	if err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}
