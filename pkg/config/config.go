package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"stori_challenge/pkg/models"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB returns the singleton instance of the database connection
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		maxRetries := 5 // Maximum number of connection attempts
		for i := 0; i < maxRetries; i++ {
			db, err = connectDB() // Attempt to connect to the database
			if err == nil {
				break // Exit the loop if connection is successful
			}
			log.Printf("Failed to connect to the database. Retrying in 5 seconds... (Attempt %d/%d)", i+1, maxRetries)
			time.Sleep(5 * time.Second) // Wait before retrying
		}
		if err != nil {
			log.Fatalf("Error connecting to the database after %d attempts: %v", maxRetries, err)
		}

		// Automatically migrate the schema to keep the database in sync with the models
		err = db.AutoMigrate(&models.SQLDocument{})
		if err != nil {
			log.Fatalf("Error migrating schema: %v", err)
		}

		log.Println("Database connection established and schema migrated")
	})

	return db // Return the database connection instance
}

// connectDB establishes a new connection to the MySQL database
func connectDB() (*gorm.DB, error) {
	// Path to the .env file at the project root
	envPath := filepath.Join(".", ".env")

	// Load the environment variables from the .env file
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	// Build the connection string using environment variables
	dbParams := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
		os.Getenv("MYSQL_CHARSET"),
		os.Getenv("MYSQL_PARSETIME"),
		os.Getenv("MYSQL_LOC"))

	// Open a new database connection
	db, err := gorm.Open(mysql.Open(dbParams), &gorm.Config{
		PrepareStmt: true, // Enable prepared statement reuse
	})
	if err != nil {
		return nil, err // Return the error if connection fails
	}

	// Configure the connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)               // Set maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)              // Set maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour * 1) // Set maximum connection lifetime

	return db, nil // Return the database connection
}
