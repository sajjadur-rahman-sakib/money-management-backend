package services

import (
	"errors"
	"money/config"
	"money/models"
	"money/utils"

	"golang.org/x/crypto/bcrypt"
)

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
