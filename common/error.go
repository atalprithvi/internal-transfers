package common

import "fmt"

type AccountError struct {
	Message string
}

func (accountError *AccountError) Error() string {
	return fmt.Sprintf("AccountError: %s", accountError.Message)
}
