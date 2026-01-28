package routes

import (
	"money/controllers"
	"money/middleware"
	"net/http"
)

func SetupRoutes() {
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()

	http.HandleFunc("/signup", authController.Signup)
	http.HandleFunc("/resend-otp", authController.ResendOTP)
	http.HandleFunc("/verify-otp", authController.VerifyOTP)
	http.HandleFunc("/login", authController.Login)
	http.HandleFunc("/profile", middleware.AuthMiddleware(userController.GetProfile))
}
