package routes

import (
	"money/controllers"
	"money/middleware"
	"net/http"
)

func SetupRoutes() {
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	bookController := controllers.NewBookController()
	transactionController := controllers.NewTransactionController()

	http.HandleFunc("/user-signup", authController.Signup)
	http.HandleFunc("/resend-otp", authController.ResendOTP)
	http.HandleFunc("/verify-otp", authController.VerifyOTP)
	http.HandleFunc("/user-login", authController.Login)
	http.HandleFunc("/user-profile", middleware.AuthMiddleware(userController.GetProfile))
	http.HandleFunc("/create-book", middleware.AuthMiddleware(bookController.CreateBook))
	http.HandleFunc("/get-books", middleware.AuthMiddleware(bookController.GetBooks))
	http.HandleFunc("/create-transaction", middleware.AuthMiddleware(transactionController.CreateTransaction))
	http.HandleFunc("/book-details", middleware.AuthMiddleware(transactionController.GetBookDetails))
}
