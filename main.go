package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {

	filePath := "txns.csv"
	err := processCSVFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
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

	for idx, row := range rows {
		// Verificar que el archivo, tenga exactamente 3 columnas
		if len(row) != 3 {
			return fmt.Errorf("la fila %d no tiene exactamente 3 columnas: %v", idx+2, row) // +2 porque la primera fila es la cabecera
		}
		fmt.Printf("%v, %v, %v\n", row[0], row[1], row[2])
	}

	return nil
}
