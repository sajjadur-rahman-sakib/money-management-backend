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
