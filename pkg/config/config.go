package config

import (
	"fmt"
	"log"
	"os"
	"stori_challenge/pkg/models"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB is the Singleton function that returns the database connection
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		maxRetries := 5
		for i := 0; i < maxRetries; i++ {
			db, err = connectDB()
			if err == nil {
				break
			}
			log.Printf("Failed to connect to database. Retrying in 5 seconds... (Attempt %d/%d)", i+1, maxRetries)
			time.Sleep(5 * time.Second)
		}
		if err != nil {
			log.Fatalf("Error connecting to database after %d attempts: %v", maxRetries, err)
		}

		// Auto migrate the schema
		err = db.AutoMigrate(&models.SQLDocument{})
		if err != nil {
			log.Fatalf("Error migrating schema: %v", err)
		}

		log.Println("Database connection established and schema migrated")
	})

	return db
}

func connectDB() (*gorm.DB, error) {
	dbParams := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
		os.Getenv("MYSQL_CHARSET"),
		os.Getenv("MYSQL_PARSETIME"),
		os.Getenv("MYSQL_LOC"))

	db, err := gorm.Open(mysql.Open(dbParams), &gorm.Config{
		PrepareStmt: true, // Reuse prepared statements
	})
	if err != nil {
		return nil, err
	}

	// Configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour * 1)

	return db, nil
}
