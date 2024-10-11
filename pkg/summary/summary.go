package summary

import (
	"fmt"
	"log"
	"stori_challenge/pkg/config"
	"stori_challenge/pkg/models"
)

// SummaryProvider defines the methods required for generating a financial summary.
type (
	SummaryProvider interface {
		TotalBalance() (float64, error)                                   // Method to retrieve the total balance
		AverageDebitAmount() (float64, error)                             // Method to retrieve the average debit amount
		AverageCreditAmount() (float64, error)                            // Method to retrieve the average credit amount
		NumberTransactionsInMonth() ([]models.TransactionsByMonth, error) // Method to retrieve transactions by month
	}

	FinanceService struct{} // Empty struct used for the FinanceService implementation
)

// CreateSummary generates a financial summary based on the provided data.
func CreateSummary(provider SummaryProvider) (models.EmailData, error) {
	// Retrieve the total balance and handle potential errors
	total, err := provider.TotalBalance()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating total balance: %w", err)
	}
	log.Printf("The total balance is: %.2f", total) // Log the total balance

	// Retrieve the average debit amount and handle potential errors
	avgDebit, err := provider.AverageDebitAmount()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating average debit amount: %w", err)
	}
	log.Printf("The average debit amount is: %.2f", avgDebit) // Log the average debit amount

	// Retrieve the average credit amount and handle potential errors
	avgCredit, err := provider.AverageCreditAmount()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating average credit amount: %w", err)
	}
	log.Printf("The average credit amount is: %.2f", avgCredit) // Log the average credit amount

	// Retrieve the number of transactions in each month and handle potential errors
	transactions, err := provider.NumberTransactionsInMonth()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error retrieving number of transactions in month: %w", err)
	}
	log.Printf("The transactions in the month are: %v", transactions) // Log the transactions by month

	// Return the compiled summary data
	return models.EmailData{
		TotalBalance:        total,
		AverageDebitAmount:  avgDebit,
		AverageCreditAmount: avgCredit,
		Transactions:        transactions,
	}, nil
}

// TotalBalance calculates the total balance from the SQLDocument model.
func (f *FinanceService) TotalBalance() (float64, error) {
	var total float64
	// Query to sum all transactions in the SQLDocument model
	if err := config.GetDB().Model(&models.SQLDocument{}).Select("SUM(transaction)").Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to get total transaction: %w", err)
	}
	return total, nil // Return the total balance
}

// AverageDebitAmount calculates the average debit amount from the SQLDocument model.
func (f *FinanceService) AverageDebitAmount() (float64, error) {
	var avg float64
	// Query to calculate the average of debit transactions (where transaction < 0)
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("transaction < ?", 0).Select("AVG(transaction)").Scan(&avg).Error; err != nil {
		return 0, fmt.Errorf("failed to get average debit transaction: %w", err)
	}
	return avg, nil // Return the average debit amount
}

// AverageCreditAmount calculates the average credit amount from the SQLDocument model.
func (f *FinanceService) AverageCreditAmount() (float64, error) {
	var avg float64
	// Query to calculate the average of credit transactions (where transaction > 0)
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("transaction > ?", 0).Select("AVG(transaction)").Scan(&avg).Error; err != nil {
		return 0, fmt.Errorf("failed to get average credit transaction: %w", err)
	}
	return avg, nil // Return the average credit amount
}

// NumberTransactionsInMonth retrieves the number of transactions for each month.
func (f *FinanceService) NumberTransactionsInMonth() ([]models.TransactionsByMonth, error) {
	transactions := []models.TransactionsByMonth{} // Slice to hold transactions by month

	// Map of month numbers to their names
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

	// Iterate through each month to count transactions
	for monthNumber, monthName := range months {
		count, err := countTransactionsByMonth(monthNumber) // Call function to count transactions for the month
		if err != nil {
			return nil, fmt.Errorf("error counting transactions for month %d: %w", monthNumber, err)
		}

		// If there are transactions for the month, append to the slice
		if count != 0 {
			newTransaction := models.TransactionsByMonth{
				Total: count,
				Month: monthName,
			}
			transactions = append(transactions, newTransaction) // Append the new transaction to the array
		}
	}

	return transactions, nil // Return the slice of transactions by month
}

// countTransactionsByMonth counts the number of transactions for a given month.
func countTransactionsByMonth(monthNumber int) (int64, error) {
	var count int64
	// Query to count transactions for the specified month
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("MONTH(date) = ?", monthNumber).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count transactions for month %d: %w", monthNumber, err)
	}
	return count, nil // Return the count of transactions
}
