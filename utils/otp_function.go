package utils

import (
	"fmt"
	"math/rand"
	"money/config"
	"net/smtp"
	"strconv"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(otp)
}

func SendOTP(email, otp string) error {
	configuration := config.GetConfig()

	smtpHost := configuration.SmtpHost
	smtpPort := configuration.SmtpPort
	smtpUsername := configuration.SmtpUsername
	smtpPassword := configuration.SmtpPassword
	smtpEmail := configuration.SmtpEmail

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	to := []string{email}
	subject := "Subject: Your OTP Code\r\n"
	body := fmt.Sprintf("Your OTP code is: %s\r\n", otp)
	msg := []byte(subject + "\r\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpEmail, to, msg)
	if err != nil {
		return err
	}

	return nil
}
