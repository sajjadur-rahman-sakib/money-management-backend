package config

import (
	"fmt"
	"log"
	"money/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDataSource() string {
	configuration := GetConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		configuration.DatabaseHost,
		configuration.DatabaseUsername,
		configuration.DatabasePassword,
		configuration.DatabaseName,
		configuration.DatabasePort,
		configuration.SslMode,
	)
	return dsn
}

func ConnectDatabase() {
	dsn := GetDataSource()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&models.User{}, &models.Book{}, &models.Transaction{})
	DB = db
}
