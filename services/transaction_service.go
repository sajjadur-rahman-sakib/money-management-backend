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
	if err := config.DB.Where("book_id = ?", bookID).Order("created_at desc").Find(&transactions).Error; err != nil {
		return nil, nil, err
	}

	return &book, transactions, nil
}

func UpdateTransaction(userID, transactionID string, amount float64, description string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := config.DB.First(&transaction, "id = ?", transactionID).Error; err != nil {
		return nil, errors.New("transaction not found")
	}

	var book models.Book
	if err := config.DB.Where("id = ? AND user_id = ?", transaction.BookID, userID).First(&book).Error; err != nil {
		return nil, errors.New("transaction not found or unauthorized")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	delta := amount - transaction.Amount
	newBalance := book.Balance
	switch transaction.Type {
	case "cash_in":
		newBalance += delta
	case "cash_out":
		newBalance -= delta
	}

	transaction.Amount = amount
	transaction.Description = description
	book.Balance = newBalance

	tx := config.DB.Begin()
	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return &transaction, nil
}

func DeleteTransaction(userID, transactionID string) error {
	var transaction models.Transaction
	if err := config.DB.First(&transaction, "id = ?", transactionID).Error; err != nil {
		return errors.New("transaction not found")
	}

	var book models.Book
	if err := config.DB.Where("id = ? AND user_id = ?", transaction.BookID, userID).First(&book).Error; err != nil {
		return errors.New("transaction not found or unauthorized")
	}

	switch transaction.Type {
	case "cash_in":
		book.Balance -= transaction.Amount
	case "cash_out":
		book.Balance += transaction.Amount
	}

	tx := config.DB.Begin()
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
