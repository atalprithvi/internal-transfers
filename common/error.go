package common

import "fmt"

type AccountError struct {
	Message string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("AccountError: %s", e.Message)
}
