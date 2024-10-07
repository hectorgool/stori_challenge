package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB es la función Singleton que devuelve la conexión a la base de datos
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		dbParams := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"),
			os.Getenv("MYSQL_CHARSET"),
			os.Getenv("MYSQL_PARSETIME"),
			os.Getenv("MYSQL_LOC"))

		db, err = gorm.Open(mysql.Open(dbParams), &gorm.Config{
			PrepareStmt: true, // Reutiliza sentencias preparadas
		})
		if err != nil {
			log.Fatalf("Error al conectar a la base de datos: %v", err)
		}

		// Configura el pool de conexiones
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Error al obtener el objeto DB: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)               // Conexiones inactivas
		sqlDB.SetMaxOpenConns(100)              // Conexiones máximas abiertas
		sqlDB.SetConnMaxLifetime(time.Hour * 1) // Duración máxima de las conexiones

		fmt.Println("Conexión a la base de datos establecida")
	})

	return db
}
