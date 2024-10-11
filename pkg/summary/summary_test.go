package summary

import (
	"testing"

	"stori_challenge/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSummaryProvider is a mock implementation of the SummaryProvider interface for testing.
type MockSummaryProvider struct {
	mock.Mock // Embedding the mock package to enable mocking behavior
}

// TotalBalance returns the total balance for the mock provider.
func (m *MockSummaryProvider) TotalBalance() (float64, error) {
	args := m.Called()                          // Call the mock's Called method
	return args.Get(0).(float64), args.Error(1) // Return the first argument and the error
}

// AverageDebitAmount returns the average debit amount for the mock provider.
func (m *MockSummaryProvider) AverageDebitAmount() (float64, error) {
	args := m.Called()                          // Call the mock's Called method
	return args.Get(0).(float64), args.Error(1) // Return the first argument and the error
}

// AverageCreditAmount returns the average credit amount for the mock provider.
func (m *MockSummaryProvider) AverageCreditAmount() (float64, error) {
	args := m.Called()                          // Call the mock's Called method
	return args.Get(0).(float64), args.Error(1) // Return the first argument and the error
}

// NumberTransactionsInMonth returns a slice of transactions aggregated by month for the mock provider.
func (m *MockSummaryProvider) NumberTransactionsInMonth() ([]models.TransactionsByMonth, error) {
	args := m.Called()                                               // Call the mock's Called method
	return args.Get(0).([]models.TransactionsByMonth), args.Error(1) // Return the first argument and the error
}

// TestCreateSummary tests the CreateSummary function using a mocked SummaryProvider.
func TestCreateSummary(t *testing.T) {
	mockProvider := new(MockSummaryProvider) // Create a new instance of the mock provider

	// Define the expected behavior for the mock methods
	mockProvider.On("TotalBalance").Return(1500.50, nil)
	mockProvider.On("AverageDebitAmount").Return(500.25, nil)
	mockProvider.On("AverageCreditAmount").Return(1000.75, nil)
	mockProvider.On("NumberTransactionsInMonth").Return([]models.TransactionsByMonth{
		{Month: "January", Total: 5},  // January transactions
		{Month: "February", Total: 3}, // February transactions
	}, nil)

	// Prepare the expected EmailData result
	expectedEmailData := models.EmailData{
		TotalBalance:        1500.50,
		AverageDebitAmount:  500.25,
		AverageCreditAmount: 1000.75,
		Transactions: []models.TransactionsByMonth{
			{Month: "January", Total: 5},
			{Month: "February", Total: 3},
		},
	}

	// Call the CreateSummary function with the mock provider
	result, err := CreateSummary(mockProvider)

	// Assert that there was no error and the result matches the expected data
	assert.NoError(t, err)                     // Check that the error is nil
	assert.Equal(t, expectedEmailData, result) // Check that the result matches the expected data

	// Verify that the mock expectations were met
	mockProvider.AssertExpectations(t) // Ensure all mocked methods were called as expected
}
