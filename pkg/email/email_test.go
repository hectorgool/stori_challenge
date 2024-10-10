package email

import (
	"stori_challenge/pkg/models"
	"testing"
)

type testpair struct {
	data   models.EmailData
	hasErr bool
}

var tests = []testpair{
	{
		models.EmailData{
			EmailTo:             "hector.gonzalez.olmos@gmail.com",
			TotalBalance:        100.0,
			AverageDebitAmount:  50.0,
			AverageCreditAmount: 150.0,
			Transactions: []models.TransactionsByMonth{
				{Total: 5, Month: "January"},
				{Total: 10, Month: "February"},
			},
		}, false,
	},
	{
		models.EmailData{
			EmailTo:             "invalid-email",
			TotalBalance:        100.0,
			AverageDebitAmount:  50.0,
			AverageCreditAmount: 150.0,
			Transactions: []models.TransactionsByMonth{
				{Total: 5, Month: "January"},
				{Total: 10, Month: "February"},
			},
		}, true,
	},
}

func TestSendEmail(t *testing.T) {
	for _, pair := range tests {
		err := SendEmail(pair.data)
		if (err != nil) != pair.hasErr {
			t.Errorf("For %+v expected error: %v, got error: %v", pair.data, pair.hasErr, err)
		}
	}
}
