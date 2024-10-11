package email

import (
	"stori_challenge/pkg/models"
	"testing"
)

// testPair defines a structure for holding test case information,
// including the email data to be tested and whether an error is expected.
type testPair struct {
	data   models.EmailData // Email data to be used in the test
	hasErr bool             // Indicates if an error is expected for this test case
}

// List of test cases with corresponding expected outcomes
var tests = []testPair{
	{
		models.EmailData{
			EmailTo:             "hector.gonzalez.olmos@gmail.com", // Valid email address
			TotalBalance:        100.0,                             // Total balance
			AverageDebitAmount:  50.0,                              // Average debit amount
			AverageCreditAmount: 150.0,                             // Average credit amount
			Transactions: []models.TransactionsByMonth{
				{Total: 5, Month: "January"},   // Transactions for January
				{Total: 10, Month: "February"}, // Transactions for February
			},
		}, false, // Expecting no error for this case
	},
	{
		models.EmailData{
			EmailTo:             "invalid-email", // Invalid email address
			TotalBalance:        100.0,           // Total balance
			AverageDebitAmount:  50.0,            // Average debit amount
			AverageCreditAmount: 150.0,           // Average credit amount
			Transactions: []models.TransactionsByMonth{
				{Total: 5, Month: "January"},   // Transactions for January
				{Total: 10, Month: "February"}, // Transactions for February
			},
		}, true, // Expecting an error for this case
	},
}

// TestSendEmail tests the SendEmail function with various inputs
func TestSendEmail(t *testing.T) {
	for _, pair := range tests {
		// Call the SendEmail function with the current test case data
		err := SendEmail(pair.data)

		// Check if the error result matches the expected outcome
		if (err != nil) != pair.hasErr {
			// Log an error if the actual result does not match the expected result
			t.Errorf("For %+v expected error: %v, got error: %v", pair.data, pair.hasErr, err)
		}
	}
}
