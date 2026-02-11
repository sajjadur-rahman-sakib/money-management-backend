package services

import (
	"errors"
	"money/config"
	"money/models"
	"money/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Signup(name, email, password, picture string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	otp := utils.GenerateOTP()
	user := &models.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Picture:   picture,
		OTP:       otp,
		OTPExpiry: time.Now().Add(5 * time.Minute),
	}

	if err := config.DB.Create(user).Error; err != nil {
		return err
	}

	if err := utils.SendOTP(email, otp); err != nil {
		return err
	}
	return nil
}

func ResendOTP(email string) error {
	var user models.User
	if err := config.DB.Where("email = ? AND is_active = false", email).First(&user).Error; err != nil {
		return errors.New("user not found or already active")
	}

	otp := utils.GenerateOTP()
	user.OTP = otp
	user.OTPExpiry = time.Now().Add(5 * time.Minute)

	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}

	if err := utils.SendOTP(email, otp); err != nil {
		return err
	}
	return nil
}

func VerifyOTP(email, otp string) error {
	var user models.User
	if err := config.DB.Where("email = ? AND otp = ?", email, otp).First(&user).Error; err != nil {
		return errors.New("invalid OTP")
	}

	if time.Now().After(user.OTPExpiry) {
		return errors.New("OTP expired")
	}

	user.IsActive = true
	user.OTP = ""
	user.OTPExpiry = time.Time{}
	return config.DB.Save(&user).Error
}

func Login(email, password string) (*models.User, string, error) {
	var user models.User

	if err := config.DB.Where("email = ? AND is_active = true", email).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func ForgotPassword(email string) error {
	var user models.User
	if err := config.DB.Where("email = ? AND is_active = true", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	otp := utils.GenerateOTP()
	user.OTP = otp
	user.OTPExpiry = time.Now().Add(5 * time.Minute)

	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}

	if err := utils.SendOTP(email, otp); err != nil {
		return err
	}
	return nil
}

func ForgotOtp(email, otp string) error {
	var user models.User
	if err := config.DB.Where("email = ? AND otp = ? AND is_active = true", email, otp).First(&user).Error; err != nil {
		return errors.New("invalid OTP")
	}

	if time.Now().After(user.OTPExpiry) {
		return errors.New("OTP expired")
	}

	return nil
}

func ResetPassword(email, otp, newPassword string) error {
	var user models.User
	if err := config.DB.Where("email = ? AND otp = ? AND is_active = true", email, otp).First(&user).Error; err != nil {
		return errors.New("invalid OTP")
	}

	if time.Now().After(user.OTPExpiry) {
		return errors.New("OTP expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.OTP = ""
	user.OTPExpiry = time.Time{}

	return config.DB.Save(&user).Error
}
