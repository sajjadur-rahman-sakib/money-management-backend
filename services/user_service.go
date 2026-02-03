package services

import (
	"errors"
	"money/config"
	"money/models"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

func GetProfile(userID string) (*models.User, error) {
	var user models.User

	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func ChangePassword(userID, currentPassword, newPassword string) error {
	var user models.User

	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return config.DB.Save(&user).Error
}

func UpdateProfile(userID, name, picturePath string) (*models.User, error) {
	var user models.User

	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if picturePath != "" {
		user.Picture = picturePath
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
