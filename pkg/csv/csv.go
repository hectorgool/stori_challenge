package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"stori_challenge/pkg/config"
	"stori_challenge/pkg/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ProcessCSVFile processes the given CSV file and stores the data in the database.
func ProcessCSVFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	if err := validateCSVHeader(reader); err != nil {
		return err
	}

	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error al leer las filas: %v", err)
	}

	return processCSVRows(rows)
}

// validateCSVHeader validates the header of the CSV file.
func validateCSVHeader(reader *csv.Reader) error {
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error al leer la cabecera: %v", err)
	}

	expectedHeaders := []string{"Id", "Date", "Transaction"}
	for i, header := range expectedHeaders {
		if strings.TrimSpace(headers[i]) != header {
			return fmt.Errorf("cabecera inválida: se esperaba %s, pero se encontró %s", header, headers[i])
		}
	}
	return nil
}

// processCSVRows processes each row in the CSV and stores them in the database.
func processCSVRows(rows [][]string) error {
	for idx, row := range rows {
		if err := validateCSVRow(row, idx); err != nil {
			return err
		}

		csvRow := models.CSVDocument{
			Id:          row[0],
			Date:        row[1],
			Transaction: row[2],
		}

		sqlDoc, err := dataCSVToSQL(csvRow)
		if err != nil {
			log.Println("Error converting CSV to SQL:", err)
			continue
		}

		if err := addTransactionToDB(sqlDoc); err != nil {
			log.Println("Error adding transaction to DB:", err)
		}
	}

	return nil
}

// validateCSVRow validates the individual row of the CSV.
func validateCSVRow(row []string, rowIndex int) error {
	if len(row) != 3 {
		return fmt.Errorf("la fila %d no tiene exactamente 3 columnas: %v", rowIndex+2, row)
	}
	return nil
}

// addTransactionToDB adds a SQLDocument to the database if it doesn't already exist.
func addTransactionToDB(sqlDoc models.SQLDocument) error {
	if err := transactionExists(sqlDoc.IdTransaction); err != nil {
		log.Println("Error:", err)
		return err
	}

	if err := config.GetDB().Create(&sqlDoc).Error; err != nil {
		return fmt.Errorf("error al crear la transacción: %v", err)
	}
	log.Println("Transacción creada exitosamente.")
	return nil
}

// transactionExists checks if a transaction already exists in the database.
func transactionExists(idTransaction uint) error {
	var existingTransaction models.SQLDocument

	if err := config.GetDB().Where("id_transaction = ?", idTransaction).First(&existingTransaction).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil // No existe, retorna nil
		}
		return err
	}

	return fmt.Errorf("la transacción con IdTransaction %d ya existe", idTransaction)
}

// dataCSVToSQL converts a CSVDocument to a SQLDocument.
func dataCSVToSQL(csvRow models.CSVDocument) (models.SQLDocument, error) {
	dateParts := strings.Split(csvRow.Date, "/")
	if len(dateParts) != 2 {
		return models.SQLDocument{}, fmt.Errorf("invalid date format for row: %v", csvRow)
	}

	month, day := dateParts[0], dateParts[1]

	IdValue, err := stringToUint(csvRow.Id)
	if err != nil {
		return models.SQLDocument{}, fmt.Errorf("error converting Id: %v", err)
	}

	TransactionFloat64, err := stringToFloat64(csvRow.Transaction)
	if err != nil {
		return models.SQLDocument{}, fmt.Errorf("error converting Transaction: %v", err)
	}

	year := getCurrentYear()
	sqlDoc := models.SQLDocument{
		IdTransaction: IdValue,
		Date:          fmt.Sprintf("%v-%v-%v", year, month, day),
		Transaction:   TransactionFloat64,
	}

	return sqlDoc, nil
}

// stringToUint converts a string to uint.
func stringToUint(s string) (uint, error) {
	num, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("error converting string to uint: %v", err)
	}
	return uint(num), nil
}

// stringToFloat64 converts a string to float64.
func stringToFloat64(s string) (float64, error) {
	floatValue, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting string to float64: %v", err)
	}
	return floatValue, nil
}

// getCurrentYear returns the current year.
func getCurrentYear() int {
	return time.Now().Year()
}

// CheckFileSize verifies if the file size is less than the specified limit in megabytes.
func CheckFileSize(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("no se pudo obtener información del archivo: %v", err)
	}

	fileSize := fileInfo.Size()
	limitBytes, err := getFileSizeLimit()
	if err != nil {
		return err
	}

	if fileSize > limitBytes {
		return fmt.Errorf("el tamaño del archivo %s (%d bytes) excede el límite de %d bytes", filePath, fileSize, limitBytes)
	}
	return nil
}

// getFileSizeLimit retrieves the file size limit from environment variable.
func getFileSizeLimit() (int64, error) {
	limitStr := os.Getenv("FILE_SIZE_LIMIT")
	if limitStr == "" {
		limitStr = "1" // 1 MB por defecto
	}

	limitMB, err := strconv.ParseFloat(limitStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error al convertir el límite de tamaño de archivo: %v", err)
	}
	return int64(limitMB * 1024 * 1024), nil // Conversión de MB a bytes
}
