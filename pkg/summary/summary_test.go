package summary

import (
	"testing"

	"stori_challenge/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSummaryProvider es un mock del SummaryProvider
type MockSummaryProvider struct {
	mock.Mock
}

func (m *MockSummaryProvider) TotalBalance() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockSummaryProvider) AverageDebitAmount() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockSummaryProvider) AverageCreditAmount() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockSummaryProvider) NumberTransactionsInMonth() ([]models.TransactionsByMonth, error) {
	args := m.Called()
	return args.Get(0).([]models.TransactionsByMonth), args.Error(1)
}

func TestCreateSummary(t *testing.T) {
	mockProvider := new(MockSummaryProvider)

	mockProvider.On("TotalBalance").Return(1500.50, nil)
	mockProvider.On("AverageDebitAmount").Return(500.25, nil)
	mockProvider.On("AverageCreditAmount").Return(1000.75, nil)
	mockProvider.On("NumberTransactionsInMonth").Return([]models.TransactionsByMonth{
		{Month: "January", Total: 5},
		{Month: "February", Total: 3},
	}, nil)

	expectedEmailData := models.EmailData{
		TotalBalance:        1500.50,
		AverageDebitAmount:  500.25,
		AverageCreditAmount: 1000.75,
		Transactions: []models.TransactionsByMonth{
			{Month: "January", Total: 5},
			{Month: "February", Total: 3},
		},
	}

	result, err := CreateSummary(mockProvider)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmailData, result)
	mockProvider.AssertExpectations(t)
}
