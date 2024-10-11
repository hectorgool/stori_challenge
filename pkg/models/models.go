package models

type (
	// CSVDocument represents the structure of a CSV file entry with fields for ID, Date, and Transaction.
	CSVDocument struct {
		Id, Date, Transaction string // Fields for ID, transaction date, and transaction details
	}

	// SQLDocument represents the structure of a SQL database entry with fields for primary key and transaction details.
	SQLDocument struct {
		Id            uint    `gorm:"primaryKey"`    // Primary key for the SQL document
		IdTransaction uint    `json:"idTransaction"` // Transaction ID for referencing the original transaction
		Date          string  `gorm:"type:date"`     // Date of the transaction in a date format
		Transaction   float64 `json:"transaction"`   // Transaction amount as a float
	}

	// TransactionsByMonth holds the total number of transactions and the corresponding month.
	TransactionsByMonth struct {
		Total int64  // Total number of transactions for the month
		Month string // Name of the month
	}

	// EmailData holds the information required for sending an email report.
	EmailData struct {
		EmailTo             string                // Recipient's email address
		TotalBalance        float64               // Total balance amount
		AverageDebitAmount  float64               // Average amount of debit transactions
		AverageCreditAmount float64               // Average amount of credit transactions
		Transactions        []TransactionsByMonth // List of transactions aggregated by month
	}
)
