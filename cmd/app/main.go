package main

import (
	"log"
	"os"
	"stori_challenge/pkg/csv"
	"stori_challenge/pkg/email"
	"stori_challenge/pkg/summary"
)

func main() {

	// Verifica que se hayan pasado exactamente dos argumentos adicionales
	if len(os.Args) != 3 {
		log.Fatal("Se requieren dos argumentos: <email> <path al archivo CSV>")
	}

	// Obtiene los argumentos de la consola
	emailWithSummary := os.Args[1]
	filePath := os.Args[2]

	// Valida el formato del correo electrónico
	if !email.IsValidEmail(emailWithSummary) {
		log.Fatalf("El formato del correo electrónico: %v es inválido", emailWithSummary)
	}

	if err := csv.CheckFileSize(filePath); err != nil {
		log.Fatalf("Error: The file %s exceeds the size limit of %s MB defined by FILE_SIZE_LIMIT. Details: %v", filePath, os.Getenv("FILE_SIZE_LIMIT"), err)
	}

	if err := csv.ProcessCSVFile(filePath); err != nil {
		log.Fatalf("Error processing CSV file (%s): %v", filePath, err)
	}

	emailData, err := summary.CreateSummary() // Ahora recibe emailData y error
	if err != nil {
		log.Fatalf("Failed to create summary: %v", err)
	}
	emailData.EmailTo = emailWithSummary
	log.Printf("Email data generated: %+v \n", emailData) // Loguear emailData

	if err := email.SendEmail(emailData); err != nil {
		log.Fatalf("Error al enviar el correo: %v\n", err)
	}

	log.Println("CSV file processed and summary created successfully.")

}
