package controller

import (
	"encoding/json"
	"fmt"
	"internal-transfers/model"
	"internal-transfers/service"
	"net/http"
)

type TransactionController struct {
	Service *service.TransactionService
}

// Handles the creation of a transaction
func (transactionController *TransactionController) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var request model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	transaction := model.Transaction{
		SourceAccountID:      request.SourceAccountID,
		DestinationAccountID: request.DestinationAccountID,
		Amount:               request.Amount,
	}

	if err := transactionController.Service.PerformTransaction(transaction); err != nil {
		http.Error(w, fmt.Sprintf("Error processing transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction successful"))
}
