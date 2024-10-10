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

	// Ruta al archivo .env en la raíz del proyecto
	envPath := filepath.Join(".", ".env")

	// Cargar el archivo .env
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error al cargar archivo: %v", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hola, Stori",
		})
	})

	r.POST("/csv", handlers.HandleCSVUpload)
	r.Run(":" + os.Getenv("HOST_PORT"))
}
