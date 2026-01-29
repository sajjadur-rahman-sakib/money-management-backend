package services

import (
	"errors"
	"money/config"
	"money/models"

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
