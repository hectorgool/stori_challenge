package summary

import (
	"fmt"
	"log"
	"stori_challenge/pkg/config"
	"stori_challenge/pkg/models"
)

// SummaryProvider defines the methods required for generating a summary.
type (
	SummaryProvider interface {
		TotalBalance() (float64, error)
		AverageDebitAmount() (float64, error)
		AverageCreditAmount() (float64, error)
		NumberTransactionsInMonth() ([]models.TransactionsByMonth, error)
	}

	FinanceService struct{}
)

// CreateSummary generates a financial summary based on the provided data.
func CreateSummary(provider SummaryProvider) (models.EmailData, error) {

	total, err := provider.TotalBalance()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating total balance: %w", err)
	}
	log.Printf("The total balance is: %.2f", total)

	avgDebit, err := provider.AverageDebitAmount()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating average debit amount: %w", err)
	}
	log.Printf("The average debit amount is: %.2f", avgDebit)

	avgCredit, err := provider.AverageCreditAmount()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating average credit amount: %w", err)
	}
	log.Printf("The average credit amount is: %.2f", avgCredit)

	transactions, err := provider.NumberTransactionsInMonth()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error retrieving number of transactions in month: %w", err)
	}
	log.Printf("The transactions in the month are: %v", transactions)

	return models.EmailData{
		TotalBalance:        total,
		AverageDebitAmount:  avgDebit,
		AverageCreditAmount: avgCredit,
		Transactions:        transactions,
	}, nil
}

// func totalBalance() (float64, error) {
func (f *FinanceService) TotalBalance() (float64, error) {
	var total float64
	if err := config.GetDB().Model(&models.SQLDocument{}).Select("SUM(transaction)").Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to get total transaction: %w", err)
	}
	return total, nil
}

// func averageDebitAmount() (float64, error) {
func (f *FinanceService) AverageDebitAmount() (float64, error) {
	var avg float64
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("transaction < ?", 0).Select("AVG(transaction)").Scan(&avg).Error; err != nil {
		return 0, fmt.Errorf("failed to get average debit transaction: %w", err)
	}
	return avg, nil
}

// func averageCreditAmount() (float64, error) {
func (f *FinanceService) AverageCreditAmount() (float64, error) {
	var avg float64
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("transaction > ?", 0).Select("AVG(transaction)").Scan(&avg).Error; err != nil {
		return 0, fmt.Errorf("failed to get average credit transaction: %w", err)
	}
	return avg, nil
}

func (f *FinanceService) NumberTransactionsInMonth() ([]models.TransactionsByMonth, error) {
	transactions := []models.TransactionsByMonth{}

	months := map[int]string{
		1:  "January",
		2:  "February",
		3:  "March",
		4:  "April",
		5:  "May",
		6:  "June",
		7:  "July",
		8:  "August",
		9:  "September",
		10: "October",
		11: "November",
		12: "December",
	}

	for monthNumber, monthName := range months {
		count, err := countTransactionsByMonth(monthNumber)
		if err != nil {
			return nil, fmt.Errorf("error counting transactions for month %d: %w", monthNumber, err)
		}

		if count != 0 {
			newTransaction := models.TransactionsByMonth{
				Total: count,
				Month: monthName,
			}
			// Append the new transaction to the array
			transactions = append(transactions, newTransaction)
		}
	}

	return transactions, nil
}

func countTransactionsByMonth(monthNumber int) (int64, error) {
	var count int64
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("MONTH(date) = ?", monthNumber).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count transactions for month %d: %w", monthNumber, err)
	}
	return count, nil
}
