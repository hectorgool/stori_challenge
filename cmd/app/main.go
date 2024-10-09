package main

import (
	"log"
	"stori_challenge/pkg/csv"
	"stori_challenge/pkg/email"
	"stori_challenge/pkg/summary"
)

func main() {
	const filePath = "txns.csv" // Usar constante para la ruta del archivo

	if err := csv.ProcessCSVFile(filePath); err != nil {
		log.Fatalf("Error processing CSV file (%s): %v", filePath, err)
	}

	emailData, err := summary.CreateSummary() // Ahora recibe emailData y error
	if err != nil {
		log.Fatalf("Failed to create summary: %v", err)
	}
	log.Printf("Email data generated: %+v \n", emailData) // Loguear emailData

	if err := email.SendEmail(emailData); err != nil {
		// Manejo de errores
		log.Fatalf("Error al enviar el correo: %v\n", err)
	}

	log.Println("CSV file processed and summary created successfully.")

}
