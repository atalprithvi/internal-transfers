package unit

import (
	"bytes"
	"encoding/json"
	"internal-transfers/controller"
	"internal-transfers/model"
	"internal-transfers/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountHandler_Success(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	mockAccount := &model.Account{
		AccountID: 1,
		Balance:   decimal.NewFromInt(100),
	}
	mockService.On("GetAccountByID", 1).Return(mockAccount, nil)

	mockAuditLogger := new(mocks.MockAuditLogger)

	controller := controller.NewAccountController(mockService, mockAuditLogger)

	req, err := http.NewRequest("GET", "/accounts/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.GetAccountHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseAccount model.Account
	err = json.NewDecoder(rr.Body).Decode(&responseAccount)
	assert.NoError(t, err)
	assert.Equal(t, 1, responseAccount.AccountID)
	assert.Equal(t, decimal.NewFromInt(100), responseAccount.Balance)

	mockService.AssertExpectations(t)
}

func TestCreateAccountHandler_Success(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	mockAccount := &model.Account{
		AccountID: 1,
		Balance:   decimal.NewFromInt(100),
	}
	mockService.On("GetAccountByID", 1).Return(nil, nil)

	mockService.On("CreateAccount", *mockAccount).Return(nil)

	mockAuditLogger := new(mocks.MockAuditLogger)

	controller := controller.NewAccountController(mockService, mockAuditLogger)

	accountInput := model.CreateAccountInput{
		AccountID:      1,
		InitialBalance: decimal.NewFromInt(100),
	}
	body, _ := json.Marshal(accountInput)

	req, err := http.NewRequest("POST", "/accounts", bytes.NewReader(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.CreateAccountHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	assert.Equal(t, "Account created successfully", rr.Body.String())

	mockService.AssertExpectations(t)
}

func TestCreateAccountHandler_AccountExists(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	mockExistingAccount := &model.Account{
		AccountID: 1,
		Balance:   decimal.NewFromInt(100),
	}
	mockService.On("GetAccountByID", 1).Return(mockExistingAccount, nil)

	mockAuditLogger := new(mocks.MockAuditLogger)

	controller := controller.NewAccountController(mockService, mockAuditLogger)

	accountInput := model.CreateAccountInput{
		AccountID:      1,
		InitialBalance: decimal.NewFromInt(100),
	}
	body, _ := json.Marshal(accountInput)

	req, err := http.NewRequest("POST", "/accounts", bytes.NewReader(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.CreateAccountHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.Equal(t, "Account already exists", rr.Body.String())

	mockService.AssertExpectations(t)
}

func TestUpdateAccountHandler_Success(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	mockAccount := &model.Account{
		AccountID: 1,
		Balance:   decimal.NewFromInt(100),
	}
	mockService.On("GetAccountByID", 1).Return(mockAccount, nil)

	mockService.On("UpdateAccountBalance", 1, decimal.NewFromInt(200)).Return(nil)

	mockAuditLogger := new(mocks.MockAuditLogger)

	controller := controller.NewAccountController(mockService, mockAuditLogger)

	accountInput := model.UpdateAccountInput{
		AccountID: 1,
		Balance:   decimal.NewFromInt(200),
	}
	body, _ := json.Marshal(accountInput)

	req, err := http.NewRequest("PUT", "/accounts/1", bytes.NewReader(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.UpdateAccountHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockService.AssertExpectations(t)
}
