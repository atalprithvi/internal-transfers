package v1

import (
	"internal-transfers/controller"
	"internal-transfers/service"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// Used to parse the request body for transaction create operation
type TransactionRequest struct {
	SourceAccountID      int             `json:"source_account_id"`
	DestinationAccountID int             `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
}

// Registers routers for transactions, version v1
func RegisterTransactionRoutes(router *mux.Router, transactionService *service.TransactionService) {
	transactionController := &controller.TransactionController{
		Service: transactionService,
	}

	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/transactions", transactionController.CreateTransactionHandler).Methods("POST")
}
