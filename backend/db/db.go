package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Use Docker's container name as the host or localhost if running locally
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	// Open the database connection
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	fmt.Println("Connected to database ✅")
}
