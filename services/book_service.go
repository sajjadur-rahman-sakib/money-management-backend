package services

import (
	"money/config"
	"money/models"

	"github.com/google/uuid"
)

func CreateBook(userID, name string) (*models.Book, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	book := &models.Book{
		Name:   name,
		UserID: parsedUUID,
	}

	if err := config.DB.Create(book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func GetBooks(userID string) ([]models.Book, error) {
	var books []models.Book

	if err := config.DB.Where("user_id = ?", userID).Find(&books).Error; err != nil {
		return nil, err
	}

	return books, nil
}
