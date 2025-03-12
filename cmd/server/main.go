package main

import (
	"fmt"
	v1 "internal-transfers/api/v1"
	"internal-transfers/common"
	"internal-transfers/persistence"
	"internal-transfers/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

// Tasks:
// 1. Establishes a connection to the database.
// 2. Initializes an audit logger.
// 3. Sets up the repositories for account and transaction data.
// 4. Initializes services for account and transaction logic.
// 5. Registers the routes for account and transaction API endpoints.
// 6. Starts the HTTP server on port 8080.
// 7. Handles graceful server shutdown upon receiving a termination signal (SIGINT, SIGTERM).

func main() {
	db, err := persistence.ConnectToDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	auditLogger, err := common.NewAuditLogger()
	if err != nil {
		log.Fatalf("Could not initialize audit logger: %v", err)
	}

	accountRepo := persistence.NewAccountRepository(db)
	transactionRepo := persistence.NewTransactionRepository(db)

	accountService := service.NewAccountService(accountRepo, auditLogger)
	transactionService := service.NewTransactionService(accountRepo, transactionRepo, auditLogger)

	router := mux.NewRouter()

	v1.RegisterAccountRoutes(router, accountService, auditLogger)

	v1.RegisterTransactionRoutes(router, transactionService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Println("Server starting on :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	fmt.Println("Shutting down server gracefully...")
	if err := server.Shutdown(nil); err != nil {
		log.Fatalf("Server Shutdown failed: %v", err)
	}
	fmt.Println("Server gracefully stopped.")
}
