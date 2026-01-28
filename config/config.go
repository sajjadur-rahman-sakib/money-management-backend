package config

import (
	"os"

	"github.com/joho/godotenv"
)

var configuration Config

type Config struct {
	DatabaseHost     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
	DatabasePort     string
	GolangPort       string
	JwtSecret        string
	SmtpHost         string
	SmtpPort         string
	SmtpUsername     string
	SmtpPassword     string
	SmtpEmail        string
}

func loadConfig() {
	err := godotenv.Load("app.env")
	if err != nil {
		panic("Error loading app.env file")
	}

	configuration = Config{
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabaseUsername: os.Getenv("DATABASE_USERNAME"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		GolangPort:       os.Getenv("GOLANG_PORT"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
		SmtpHost:         os.Getenv("SMTP_HOST"),
		SmtpPort:         os.Getenv("SMTP_PORT"),
		SmtpUsername:     os.Getenv("SMTP_USERNAME"),
		SmtpPassword:     os.Getenv("SMTP_PASSWORD"),
		SmtpEmail:        os.Getenv("SMTP_EMAIL"),
	}
}

func GetConfig() Config {
	loadConfig()
	return configuration
}
