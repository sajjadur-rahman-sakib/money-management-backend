package controllers

import (
	"encoding/json"
	"money/services"
	"money/utils"
	"net/http"
)

type TransactionController struct{}

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

func (transactionController *TransactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var request struct {
		BookID      string  `json:"book_id"`
		Type        string  `json:"type"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	transaction, err := services.CreateTransaction(userID, request.BookID, request.Type, request.Amount, request.Description)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (transactionController *TransactionController) GetBookDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var request struct {
		BookID string `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if request.BookID == "" {
		utils.RespondError(w, http.StatusBadRequest, "book_id required")
		return
	}

	book, transactions, err := services.GetBookDetails(userID, request.BookID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"book":         book,
		"transactions": transactions,
	})
}

func (transactionController *TransactionController) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var request struct {
		TransactionID string  `json:"transaction_id"`
		Amount        float64 `json:"amount"`
		Description   string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if request.TransactionID == "" {
		utils.RespondError(w, http.StatusBadRequest, "transaction_id required")
		return
	}

	if request.Amount <= 0 {
		utils.RespondError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	transaction, err := services.UpdateTransaction(userID, request.TransactionID, request.Amount, request.Description)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (transactionController *TransactionController) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var request struct {
		TransactionID string `json:"transaction_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if request.TransactionID == "" {
		utils.RespondError(w, http.StatusBadRequest, "transaction_id required")
		return
	}

	if err := services.DeleteTransaction(userID, request.TransactionID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted"})
}
