package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	//.Println("Stori Challenge")
	filePath := "txns.csv"
	err := processCSVFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("El archivo CSV es válido.")
	}
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

	// Leer todas las filas
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error al leer las filas: %v", err)
	}

	for _, row := range rows {
		fmt.Printf("%v, %v, %v\n", row[0], row[1], row[2])
	}

	return nil
}
