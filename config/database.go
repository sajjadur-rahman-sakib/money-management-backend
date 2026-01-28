package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDataSource() string {
	configuration := GetConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		configuration.DatabaseHost,
		configuration.DatabaseUsername,
		configuration.DatabasePassword,
		configuration.DatabaseName,
		configuration.DatabasePort,
	)
	return dsn
}

func ConnectDatabase() {
	dsn := GetDataSource()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	//db.AutoMigrate(&models.User{})
	DB = db
}
