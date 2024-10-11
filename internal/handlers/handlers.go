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

// HandleCSVUpload handles the CSV file upload and summary creation.
func HandleCSVUpload(c *gin.Context) {
	// Retrieve the email address from the form data
	emailWithSummary := c.PostForm("email")

	// Validate the format of the email address
	if !email.IsValidEmail(emailWithSummary) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Retrieve the uploaded CSV file from the form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve the file"})
		return
	}

	// Create a temporary file to store the uploaded CSV
	tempFile, err := os.CreateTemp("", "csv-*.csv")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create temporary file"})
		return
	}
	defer os.Remove(tempFile.Name()) // Ensure the temporary file is removed after use
	defer tempFile.Close()           // Close the temporary file when done

	// Save the uploaded file to the temporary location
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save the file"})
		return
	}

	// Check the size of the uploaded file
	if err := csv.CheckFileSize(tempFile.Name()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File exceeds the allowed size limit"})
		return
	}

	// Process the uploaded CSV file
	if err := csv.ProcessCSVFile(tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing the CSV file"})
		return
	}

	// Create the summary from the processed data
	provider := &summary.FinanceService{}
	emailData, err := summary.CreateSummary(provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the summary"})
		return
	}
	emailData.EmailTo = emailWithSummary // Set the recipient email address

	// Send the summary email
	if err := email.SendEmail(emailData); err != nil {
		log.Printf("Error sending email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending the email"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "CSV file processed and summary sent successfully"})
}
