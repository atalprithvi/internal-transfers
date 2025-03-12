package v1

import (
	"internal-transfers/common"
	"internal-transfers/controller"
	"internal-transfers/service"

	"github.com/gorilla/mux"
)

// Registers routers for accounts, version v1
func RegisterAccountRoutes(router *mux.Router, accountService *service.AccountService, auditLogger *common.AuditLogger) {
	accountController := controller.NewAccountController(accountService, auditLogger)

	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/accounts/{account_id}", accountController.UpdateAccountHandler).Methods("PUT")
	router.HandleFunc("/api/v1/accounts/{account_id:[0-9]+}", accountController.GetAccountHandler).Methods("GET")
	router.HandleFunc("/api/v1/accounts/{account_id:[0-9]+}", accountController.GetAccountHandler).Methods("GET")
	router.HandleFunc("/api/v1/accounts", accountController.CreateAccountHandler).Methods("POST")
}
