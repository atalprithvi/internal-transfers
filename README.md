# Internal Transfers Application

This is a Go-based application that facilitates internal financial transfers between accounts. The application exposes HTTP endpoints to create accounts, query account balances, and submit financial transactions.

About the Project
/internal-transfers
├── /cmd
│   └── main.go                             # Entry point for the application
├── /common
│   ├── audit_logger.go                    # Implements the AuditLogger for logging application actions
│   ├── config.go                          # Configuration utilities, including loading environment variables
│   └── logger.go                          # Logger utilities for consistent logging across the application
├── /controller
│   ├── account_controller.go              # Handles HTTP requests for account-related operations (REST API controllers)
│   └── transaction_controller.go          # Handles HTTP requests for transaction-related operations
├── /model
│   ├── account.go                         # Defines the Account data structure
│   ├── transaction.go                     # Defines the Transaction data structure
│   └── input_validation.go                # Input validation logic and request data models (e.g., CreateAccountInput)
├── /persistence
│   ├── account_repository.go              # Implements data access methods for account operations (CRUD)
│   ├── transaction_repository.go          # Implements data access methods for transaction operations (CRUD)
│   └── db.go                              # Database connection and setup
├── /service
│   ├── account_service.go                 # Business logic for account operations
│   ├── transaction_service.go             # Business logic for transaction operations
├── /tests
│   ├── unit
│   │   ├── account_service_test.go        # Unit tests for account service logic
│   │   ├── transaction_service_test.go    # Unit tests for transaction service logic
│   │   ├── account_controller_test.go     # Unit tests for account controller
│   │   └── transaction_controller_test.go # Unit tests for transaction controller
│   ├── mocks
│   │   ├── mock_account_repository.go    # Mock for AccountRepository for unit tests
│   │   ├── mock_transaction_repository.go# Mock for TransactionRepository for unit tests
│   │   └── mock_account_service.go       # Mock for AccountService for unit tests
└── go.mod                                 # Go modules file for dependencies



## Prerequisites

- Go 1.18+
- PostgreSQL

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/internal-transfers.git
   cd internal-transfers
How to Run:
1. go to internal-transfers/cmd/server
2. go run .


curl -X POST http://localhost:8080/api/v1/accounts \
-H "Content-Type: application/json" \
-d '{
  "account_id": 123,
  "initial_balance": "1000.24"
}'


curl -X POST http://localhost:8080/api/v1/transactions \
-H "Content-Type: application/json" \
-d '{
  "source_account_id": 123,
  "destination_account_id": 345,
  "amount": "568.90"
}'


curl -X GET http://localhost:8080/api/v1/accounts/13


curl -X POST http://localhost:8080/api/v1/accounts \
-H "Content-Type: application/json" \
-d '{
  "account_id": 888,
  "initial_balance": "100.23344"
}'


DB Script:
CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    balance DECIMAL(15, 5) NOT NULL
);


CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    source_account_id INT NOT NULL,
    destination_account_id INT NOT NULL,
    amount DECIMAL(15, 5) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_account_id) REFERENCES accounts(account_id),
    FOREIGN KEY (destination_account_id) REFERENCES accounts(account_id)
);

*****
   Assumptions:
    - Each account must have a unique account_id.
    - An account balance is a decimal number that can hold fractional amounts.
    - The account_id is a mandatory field when creating an account.
    - Both source_account_id and destination_account_id must refer to existing, valid accounts.
    - The transaction amount must be a positive decimal number.
    - Sufficient balance should be available in the source account to perform a transaction.
    - Both the accounts are in an valid stat. Other states e.g. inactive, frozen are out of scope.
    
***
TODO:
I) Complete Unit tests coverage.
II)Add created_date/udpated_date,updated_by columns.
III) Handle concurrency.
