package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"stori_challenge/config"
	"stori_challenge/models"
	"strconv"
	"strings"
	"text/template"
	"time"

	"gorm.io/gorm"
)

func main() {
	const filePath = "txns.csv" // Usar constante para la ruta del archivo

	if err := processCSVFile(filePath); err != nil {
		log.Fatalf("Error processing CSV file (%s): %v", filePath, err)
	}

	emailData, err := createSummary() // Ahora recibe emailData y error
	if err != nil {
		log.Fatalf("Failed to create summary: %v", err)
	}
	log.Printf("Email data generated: %+v \n", emailData) // Loguear emailData

	/*
		if err := SendEmail(emailData); err != nil {
			// Manejo de errores
			log.Fatalf("Error al enviar el correo: %v\n", err)
		}
	*/

	log.Println("CSV file processed and summary created successfully.")

}

// createSummary genera un resumen de los datos financieros y maneja errores de manera más robusta.
func createSummary() (models.EmailData, error) {

	total, err := totalBalance()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating total balance: %w", err) // Propaga el error
	}
	log.Printf("The total balance is: %.2f", total)

	avgDebit, err := averageDebitAmount()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating average debit amount: %w", err)
	}
	log.Printf("The average debit amount is: %.2f", avgDebit)

	avgCredit, err := averageCreditAmount()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error calculating average credit amount: %w", err)
	}
	log.Printf("The average credit amount is: %.2f", avgCredit)

	transactions, err := numberTransactionsInMonth()
	if err != nil {
		return models.EmailData{}, fmt.Errorf("error retrieving transactions by month: %w", err)
	}

	var emailData models.EmailData
	emailData.EmailTo = "hector.gonzalez.olmos@gmail.com"
	emailData.TotalBalance = total
	emailData.AverageDebitAmount = avgDebit
	emailData.AverageCreditAmount = avgCredit
	emailData.Transactions = transactions

	return emailData, nil // Regresa emailData y nil en caso de éxito
}

func totalBalance() (float64, error) {
	var total float64
	if err := config.GetDB().Model(&models.SQLDocument{}).Select("SUM(transaction)").Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to get total transaction: %w", err)
	}
	return total, nil
}

func averageDebitAmount() (float64, error) {
	var avg float64
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("transaction < ?", 0).Select("AVG(transaction)").Scan(&avg).Error; err != nil {
		return 0, fmt.Errorf("failed to get average debit transaction: %w", err)
	}
	return avg, nil
}

func averageCreditAmount() (float64, error) {
	var avg float64
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("transaction > ?", 0).Select("AVG(transaction)").Scan(&avg).Error; err != nil {
		return 0, fmt.Errorf("failed to get average credit transaction: %w", err)
	}
	return avg, nil
}

func countTransactionsByMonth(monthNumber int) (int64, error) {
	var count int64
	if err := config.GetDB().Model(&models.SQLDocument{}).Where("MONTH(date) = ?", monthNumber).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count transactions for month %d: %w", monthNumber, err)
	}
	return count, nil
}

func numberTransactionsInMonth() ([]models.TransactionsByMonth, error) {
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

func processCSVFile(filePath string) error {
	// Abrir el archivo CSV
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Leer el archivo CSV
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Leer la cabecera
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error al leer la cabecera: %v", err)
	}

	// Verificar si la cabecera es correcta
	expectedHeaders := []string{"Id", "Date", "Transaction"}
	for i, header := range expectedHeaders {
		if strings.TrimSpace(headers[i]) != header {
			return fmt.Errorf("cabecera inválida: se esperaba %s, pero se encontró %s", header, headers[i])
		}
	}

	// Leer todas las filas
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error al leer las filas: %v", err)
	}

	var csvRow models.CSVDocument
	for idx, row := range rows {
		// Verificar que el archivo, tenga exactamente 3 columnas
		if len(row) != 3 {
			return fmt.Errorf("la fila %d no tiene exactamente 3 columnas: %v", idx+2, row) // +2 porque la primera fila es la cabecera
		}
		csvRow.Id = row[0]
		csvRow.Date = row[1]
		csvRow.Transaction = row[2]

		sqlDoc, err := dataCSVToSQL(csvRow)
		if err != nil {
			// Si ocurre un error, se imprime y se termina la ejecución
			fmt.Println("Error converting CSV to SQL:", err)
		}
		fmt.Println(sqlDoc)
		addTransactionToDB(sqlDoc)

	}

	return nil
}

func addTransactionToDB(sqlDoc models.SQLDocument) {
	// Verifica si la transacción ya existe llamando a transactionExists
	if err := transactionExists(sqlDoc.IdTransaction); err != nil {
		// Si existe o hay un error, maneja la salida
		log.Println("Error:", err)
		return
	}

	// Si no existe, procede a crear el nuevo registro
	if err := config.GetDB().Create(&sqlDoc).Error; err != nil {
		log.Fatalln("Error al crear la transacción:", err)
	} else {
		log.Println("Transacción creada exitosamente.")
	}
}

func transactionExists(idTransaction uint) error {
	var existingTransaction models.SQLDocument

	// Busca si ya existe un registro con el IdTransaction
	if err := config.GetDB().Where("id_transaction = ?", idTransaction).First(&existingTransaction).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			// Si no se encuentra un registro, no hay error, retorna nil
			return nil
		}
		// Si ocurre otro error, lo retorna
		return err
	}

	// Si encuentra un registro, retorna un error personalizado
	return fmt.Errorf("la transacción con IdTransaction %d ya existe", idTransaction)
}

// Convertir un registro CSV a un registro SQL
func dataCSVToSQL(csvRow models.CSVDocument) (models.SQLDocument, error) {
	// Validación y procesamiento del campo Date
	dateParts := strings.Split(csvRow.Date, "/")
	if len(dateParts) != 2 {
		return models.SQLDocument{}, fmt.Errorf("invalid date format for row: %v", csvRow)
	}

	month, day := dateParts[0], dateParts[1]

	// Conversión de string a uint para el campo Id
	IdValue, err := stringToUint(csvRow.Id)
	if err != nil {
		return models.SQLDocument{}, fmt.Errorf("error converting Id: %v", err)
	}

	// Conversión de string a float64 para el campo Transaction
	TransactionFloat64, err := stringToFloat64(csvRow.Transaction)
	if err != nil {
		return models.SQLDocument{}, fmt.Errorf("error converting Transaction: %v", err)
	}

	year := getCurrentYear()
	// Creación del documento SQL con los datos procesados
	sqlDoc := models.SQLDocument{
		IdTransaction: IdValue,
		Date:          fmt.Sprintf("%v-%v-%v", year, month, day), // Fecha en formato YYYY-MM-DD
		Transaction:   TransactionFloat64,
	}

	return sqlDoc, nil
}

// Convertir el Id string a uint
// Recibe un string, lo convierte a un número sin signo (uint) y maneja errores
func stringToUint(s string) (uint, error) {
	num, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("error converting string to uint: %v", err)
	}
	return uint(num), nil
}

// Convertir Transaction string a float64
// Recibe un string, lo convierte a float64 y maneja errores
func stringToFloat64(s string) (float64, error) {
	floatValue, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting string to float64: %v", err)
	}
	return floatValue, nil
}

func getCurrentYear() int {
	return time.Now().Year()
}

func SendEmail(data models.EmailData) error {
	// Configuración del servidor SMTP
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_SENDER")
	senderPassword := os.Getenv("SMTP_PASSWD")

	// Validación de variables de entorno
	if smtpServer == "" || smtpPort == "" || senderEmail == "" || senderPassword == "" {
		return fmt.Errorf("missing SMTP configuration in environment variables")
	}

	// Autenticación con el servidor SMTP
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	// Plantilla HTML externa
	templateFile := "email_template.html"

	// Parseamos la plantilla HTML
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Creamos un buffer para almacenar la salida de la plantilla
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Mensaje del correo electrónico
	htmlMessage := tpl.String()

	// Destinatarios
	to := []string{os.Getenv("SMTP_SENDER")}
	to = append(to, data.EmailTo)

	// Destinatarios en copia (CC)
	cc := os.Getenv("SMTP_CC")
	ccEmails := []string{}
	if cc != "" {
		ccEmails = append(ccEmails, cc)
	}

	subject := os.Getenv("SMTP_SUBJECT")

	// Cuerpo del correo
	body := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nCc: %s\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		strings.Join(to, ","),
		subject,
		strings.Join(ccEmails, ","),
		htmlMessage))

	// Combinar destinatarios y copia
	recipients := append(to, ccEmails...)

	// Enviar el correo
	if err = smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, recipients, body); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
