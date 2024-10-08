package models

type (
	CSVDocument struct {
		Id, Date, Transaction string
	}
	SQLDocument struct {
		Id            uint    `gorm:"primaryKey"`
		IdTransaction uint    `json:"idTransaction"`
		Date          string  `gorm:"type:date"`
		Transaction   float64 `json:"transaction"`
	}
	TransactionsByMonth struct {
		Total int64
		Month string
	}
)
