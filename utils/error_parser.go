package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeAuthentication = "AUTHENTICATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeServer         = "SERVER_ERROR"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeBadRequest     = "BAD_REQUEST"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RespondError(w http.ResponseWriter, status int, message string) {
	code := getErrorCode(status, message)
	userMessage := getUserFriendlyMessage(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIError{
		Code:    code,
		Message: userMessage,
	})
}

func getErrorCode(status int, message string) string {
	switch status {
	case http.StatusBadRequest:
		if containsAny(message, "password", "email", "name", "required", "invalid") {
			return ErrCodeValidation
		}
		return ErrCodeBadRequest
	case http.StatusUnauthorized:
		return ErrCodeAuthentication
	case http.StatusForbidden:
		return ErrCodeUnauthorized
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusConflict:
		return ErrCodeConflict
	default:
		return ErrCodeServer
	}
}

func getUserFriendlyMessage(message string) string {
	lowerMsg := strings.ToLower(message)

	if strings.Contains(lowerMsg, "duplicate key") && strings.Contains(lowerMsg, "email") {
		return "This email is already registered. Please use a different email or try logging in."
	}
	if strings.Contains(lowerMsg, "unique constraint") && strings.Contains(lowerMsg, "email") {
		return "This email is already registered. Please use a different email or try logging in."
	}
	if strings.Contains(lowerMsg, "already exists") || strings.Contains(lowerMsg, "already registered") {
		return "This email is already registered. Please use a different email or try logging in."
	}

	if strings.Contains(lowerMsg, "invalid otp") {
		return "The verification code you entered is incorrect. Please check and try again."
	}
	if strings.Contains(lowerMsg, "otp expired") {
		return "Your verification code has expired. Please request a new one."
	}

	if strings.Contains(lowerMsg, "current password is incorrect") {
		return "The current password you entered is incorrect. Please try again."
	}
	if strings.Contains(lowerMsg, "passwords do not match") {
		return "The passwords you entered don't match. Please make sure both passwords are the same."
	}
	if strings.Contains(lowerMsg, "invalid credentials") {
		return "Invalid email or password. Please check your credentials and try again."
	}

	if strings.Contains(lowerMsg, "user not found") {
		return "We couldn't find an account with that email. Please check the email or sign up."
	}
	if strings.Contains(lowerMsg, "user not found or already active") {
		return "This account is either not found or already activated. Please try logging in."
	}

	if strings.Contains(lowerMsg, "required") {
		return "Please fill in all required fields."
	}
	if strings.Contains(lowerMsg, "invalid json") {
		return "Invalid request format. Please try again."
	}

	if strings.Contains(lowerMsg, "unauthorized") || strings.Contains(lowerMsg, "token") {
		return "Your session has expired. Please log in again."
	}

	if strings.Contains(lowerMsg, "connection") || strings.Contains(lowerMsg, "timeout") {
		return "Unable to connect to the server. Please check your internet connection."
	}

	if isUserFriendly(message) {
		return message
	}
	return "Something went wrong. Please try again later."
}

func containsAny(s string, substrings ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range substrings {
		if strings.Contains(lower, sub) {
			return true
		}
	}
	return false
}

func isUserFriendly(message string) bool {
	if len(message) == 0 {
		return false
	}
	lowerMsg := strings.ToLower(message)
	technicalTerms := []string{"error:", "exception", "stack trace", "panic", "nil pointer", "sql", "gorm"}
	for _, term := range technicalTerms {
		if strings.Contains(lowerMsg, term) {
			return false
		}
	}
	return true
}
