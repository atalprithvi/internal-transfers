package mocks

import (
	"internal-transfers/model"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

// MockAccountService is a mock implementation of the AccountService interface
type MockAccountService struct {
	mock.Mock
}

// Implement the GetAccountByID method (from service.AccountService)
func (m *MockAccountService) GetAccountByID(accountID int) (*model.Account, error) {
	args := m.Called(accountID)
	return args.Get(0).(*model.Account), args.Error(1)
}

// Implement the CreateAccount method (from service.AccountService)
func (m *MockAccountService) CreateAccount(account model.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

// Implement the UpdateAccountBalance method (from service.AccountService)
func (m *MockAccountService) UpdateAccountBalance(accountID int, newBalance decimal.Decimal) error {
	args := m.Called(accountID, newBalance)
	return args.Error(0)
}
