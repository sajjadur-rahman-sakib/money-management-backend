package services

import (
	"errors"
	"money/config"
	"money/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func UpdateBookName(userID, bookID, name string) (*models.Book, error) {
	var book models.Book
	if err := config.DB.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found or unauthorized")
		}
		return nil, err
	}

	book.Name = name
	if err := config.DB.Save(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func DeleteBook(userID, bookID string) error {
	var book models.Book
	if err := config.DB.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("book not found or unauthorized")
		}
		return err
	}

	tx := config.DB.Begin()
	if err := tx.Where("book_id = ?", bookID).Delete(&models.Transaction{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&book).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
