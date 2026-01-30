package services

import (
	"errors"
	"money/config"
	"money/models"
)

func CreateTransaction(userID, bookID, transactionType string, amount float64, description string) (*models.Transaction, error) {
	var book models.Book

	if err := config.DB.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
		return nil, errors.New("book not found or unauthorized")
	}

	if transactionType != "cash_in" && transactionType != "cash_out" {
		return nil, errors.New("invalid transaction type")
	}

	if transactionType == "cash_out" && book.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	transaction := &models.Transaction{
		BookID:      book.ID,
		Type:        transactionType,
		Amount:      amount,
		Description: description,
	}

	if transactionType == "cash_in" {
		book.Balance += amount
	} else {
		book.Balance -= amount
	}

	tx := config.DB.Begin()
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return transaction, nil
}

func GetBookDetails(userID, bookID string) (*models.Book, []models.Transaction, error) {
	var book models.Book
	if err := config.DB.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
		return nil, nil, errors.New("book not found or unauthorized")
	}

	var transactions []models.Transaction
	if err := config.DB.Where("book_id = ?", bookID).Find(&transactions).Error; err != nil {
		return nil, nil, err
	}

	return &book, transactions, nil
}
