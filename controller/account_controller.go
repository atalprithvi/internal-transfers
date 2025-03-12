package controller

import (
	"encoding/json"
	"fmt"
	"internal-transfers/common"
	"internal-transfers/model"
	"internal-transfers/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// Handles the HTTP requests for account related operations
type AccountController struct {
	Service     *service.AccountService
	AuditLogger *common.AuditLogger
}

func NewAccountController(accountService *service.AccountService, auditLogger *common.AuditLogger) *AccountController {
	return &AccountController{
		Service:     accountService,
		AuditLogger: auditLogger,
	}
}

// Retrieve an account by its ID
func (ac *AccountController) GetAccountHandler(writer http.ResponseWriter, request *http.Request) {
	accountID := mux.Vars(request)["account_id"]

	id, err := strconv.Atoi(accountID)
	if err != nil {
		http.Error(writer, "Invalid account ID format", http.StatusBadRequest)
		return
	}

	account, err := ac.Service.GetAccountByID(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error fetching account: %v", err), http.StatusInternalServerError)
		return
	}

	if account == nil {
		http.Error(writer, fmt.Sprintf("Account with ID %d not found", id), http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(account); err != nil {
		http.Error(writer, fmt.Sprintf("Error encoding account data: %v", err), http.StatusInternalServerError)
	}
}

func (accountController *AccountController) CreateAccountHandler(writer http.ResponseWriter, r *http.Request) {
	var input model.CreateAccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(writer, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.AccountID == 0 {
		http.Error(writer, "Account ID must be provided", http.StatusBadRequest)
		return
	}

	existingAccount, err := accountController.Service.GetAccountByID(input.AccountID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error checking account: %v", err), http.StatusInternalServerError)
		return
	}
	if existingAccount != nil {
		http.Error(writer, "Account already exists", http.StatusBadRequest)
		return
	}

	initialBalance, err := decimal.NewFromString(input.InitialBalance.String())
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid initial balance format: %v", err), http.StatusBadRequest)
		return
	}

	newAccount := model.NewAccount(input.AccountID, initialBalance)

	if err := accountController.Service.CreateAccount(*newAccount); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	accountController.AuditLogger.LogAction("CreateAccount", fmt.Sprintf("Account created with ID: %d", input.AccountID))

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte("Account created successfully"))
}

// Handles the PUT /v1/accounts/{account_id} request
func (accountController *AccountController) UpdateAccountHandler(writer http.ResponseWriter, request *http.Request) {
	var input model.UpdateAccountInput
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(writer, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.AccountID == 0 {
		http.Error(writer, "Please provide a valid account_id", http.StatusBadRequest)
		return
	}

	existingAccount, err := accountController.Service.GetAccountByID(input.AccountID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error checking account: %v", err), http.StatusInternalServerError)
		return
	}
	if existingAccount == nil {
		http.Error(writer, "Account not found", http.StatusNotFound)
		return
	}

	if err := accountController.Service.UpdateAccountBalance(input.AccountID, input.Balance); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	accountController.AuditLogger.LogAction("UpdateAccount", fmt.Sprintf("Account updated with ID: %d", input.AccountID))

	writer.WriteHeader(http.StatusOK)
}
