package mocks

import (
	"context"
	"internal-transfers/model"

	"github.com/shopspring/decimal"
)

type MockAccountRepository struct {
	MockGetAccountByIDWithContext       func(ctx context.Context, accountID int) (*model.Account, error)
	MockCreateAccountWithContext        func(ctx context.Context, account model.Account) error
	MockUpdateAccountBalanceWithContext func(ctx context.Context, accountID int, newBalance decimal.Decimal) error
}

func (m *MockAccountRepository) GetAccountByIDWithContext(ctx context.Context, accountID int) (*model.Account, error) {
	return m.MockGetAccountByIDWithContext(ctx, accountID)
}

func (m *MockAccountRepository) CreateAccountWithContext(ctx context.Context, account model.Account) error {
	return m.MockCreateAccountWithContext(ctx, account)
}

func (m *MockAccountRepository) UpdateAccountBalanceWithContext(ctx context.Context, accountID int, newBalance decimal.Decimal) error {
	return m.MockUpdateAccountBalanceWithContext(ctx, accountID, newBalance)
}
