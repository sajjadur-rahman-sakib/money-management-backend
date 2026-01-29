package controllers

import (
	"encoding/json"
	"money/services"
	"net/http"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (authController *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := services.Login(request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"picture": user.Picture,
		"name":    user.Name,
		"email":   user.Email,
		"token":   token,
	})
}
