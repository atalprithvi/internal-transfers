package mocks

import (
	"context"
	"internal-transfers/model"
	"internal-transfers/persistence"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

var _ persistence.AccountRepository = (*MockAccountRepository)(nil)

func (m *MockAccountRepository) GetAccountByIDWithContext(ctx context.Context, accountID int) (*model.Account, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) CreateAccountWithContext(ctx context.Context, account model.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) UpdateAccountBalanceWithContext(ctx context.Context, accountID int, newBalance decimal.Decimal) error {
	args := m.Called(ctx, accountID, newBalance)
	return args.Error(0)
}
