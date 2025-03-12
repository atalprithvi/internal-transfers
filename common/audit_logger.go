package common

import (
	"fmt"
	"log"
	"os"
)

// Responsible for logging audit actions
type AuditLogger struct {
	logger *log.Logger
}

func NewAuditLogger() (*AuditLogger, error) {
	logFile, err := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}
	return &AuditLogger{
		logger: log.New(logFile, "", log.LstdFlags),
	}, nil
}

// Logs a specific action
func (a *AuditLogger) LogAction(action string, details string) {
	if a.logger != nil {
		a.logger.Printf("%s: %s\n", action, details)
	} else {
		fmt.Printf("AuditLogger is not initialized properly. Action: %s, Details: %s\n", action, details)
	}
}
