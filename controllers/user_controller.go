package controllers

import (
	"encoding/json"
	"money/services"
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (userController *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := services.GetProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"picture": user.Picture,
	})
}

func (userController *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	if err := services.ChangePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated"})
}
