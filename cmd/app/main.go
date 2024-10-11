package main

import (
	"log"
	"os"
	"path/filepath"
	"stori_challenge/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Define the path to the .env file located at the root of the project
	envPath := filepath.Join(".", ".env")

	// Load the environment variables from the .env file
	err := godotenv.Load(envPath)
	if err != nil {
		// Log a fatal error and terminate the program if loading fails
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize a new Gin router
	r := gin.Default()

	// Define a GET endpoint for the root URL that responds with a welcome message
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Stori",
		})
	})

	// Define a POST endpoint for uploading CSV files, delegating to the HandleCSVUpload handler
	r.POST("/csv", handlers.HandleCSVUpload)

	// Start the Gin server on the specified host port from environment variables
	r.Run(":" + os.Getenv("HOST_PORT"))
}
