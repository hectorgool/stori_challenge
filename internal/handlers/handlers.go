package handlers

import (
	"log"
	"net/http"
	"os"
	"stori_challenge/pkg/csv"
	"stori_challenge/pkg/email"
	"stori_challenge/pkg/summary"

	"github.com/gin-gonic/gin"
)

func HandleCSVUpload(c *gin.Context) {
	// Obtener el email del formulario
	emailWithSummary := c.PostForm("email")

	// Validar el formato del correo electrónico
	if !email.IsValidEmail(emailWithSummary) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de correo electrónico inválido"})
		return
	}

	// Obtener el archivo CSV
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo obtener el archivo"})
		return
	}

	// Crear un archivo temporal para guardar el CSV
	tempFile, err := os.CreateTemp("", "csv-*.csv")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear un archivo temporal"})
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Guardar el archivo subido en el archivo temporal
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el archivo"})
		return
	}

	// Verificar el tamaño del archivo
	if err := csv.CheckFileSize(tempFile.Name()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo excede el límite de tamaño permitido"})
		return
	}

	// Procesar el archivo CSV
	if err := csv.ProcessCSVFile(tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar el archivo CSV"})
		return
	}

	// Crear el resumen
	provider := &summary.FinanceService{}
	emailData, err := summary.CreateSummary(provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el resumen"})
		return
	}
	emailData.EmailTo = emailWithSummary

	// Enviar el email
	if err := email.SendEmail(emailData); err != nil {
		log.Printf("Error al enviar correo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el correo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Archivo CSV procesado y resumen enviado exitosamente"})
}
