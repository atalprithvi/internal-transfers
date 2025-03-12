package unit

import (
	"context"
	"internal-transfers/common"
	"internal-transfers/model"
	"internal-transfers/service"
	"internal-transfers/tests/mocks"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	mockRepo := &mocks.MockAccountRepository{
		MockGetAccountByIDWithContext: func(ctx context.Context, accountID int) (*model.Account, error) {
			return nil, nil
		},
		MockCreateAccountWithContext: func(ctx context.Context, account model.Account) error {
			return nil
		},
	}

	auditLogger := &common.AuditLogger{}
	accountService := service.NewAccountService(mockRepo, auditLogger)

	account := model.Account{
		AccountID: 1,
		Balance:   decimal.NewFromInt(100),
	}

	err := accountService.CreateAccount(account)

	assert.NoError(t, err, "Expected no error while creating the account")
}

func TestCreateAccount_AccountExists(t *testing.T) {
	mockRepo := &mocks.MockAccountRepository{
		MockGetAccountByIDWithContext: func(ctx context.Context, accountID int) (*model.Account, error) {
			return &model.Account{AccountID: accountID, Balance: decimal.NewFromInt(100)}, nil
		},
		MockCreateAccountWithContext: func(ctx context.Context, account model.Account) error {
			return nil
		},
	}

	auditLogger := &common.AuditLogger{}
	accountService := service.NewAccountService(mockRepo, auditLogger)

	account := model.Account{
		AccountID: 1,
		Balance:   decimal.NewFromInt(100),
	}

	err := accountService.CreateAccount(account)

	assert.Error(t, err, "Expected error because the account already exists")
}
