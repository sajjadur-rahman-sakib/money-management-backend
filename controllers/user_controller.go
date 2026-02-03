package controllers

import (
	"encoding/json"
	"errors"
	"money/services"
	"money/utils"
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (userController *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := services.GetProfile(userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
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
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		utils.RespondError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		utils.RespondError(w, http.StatusBadRequest, "Passwords do not match")
		return
	}

	if err := services.ChangePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated"})
}

func (userController *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	name := r.FormValue("name")

	var picturePath string
	file, header, err := r.FormFile("picture")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			utils.RespondError(w, http.StatusBadRequest, "Failed to read picture")
			return
		}
	} else {
		defer file.Close()
		picturePath, err = utils.SavePicture(file, header)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to save picture")
			return
		}
	}

	if name == "" && picturePath == "" {
		utils.RespondError(w, http.StatusBadRequest, "Nothing to update")
		return
	}

	user, err := services.UpdateProfile(userID, name, picturePath)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"picture": user.Picture,
	})
}
