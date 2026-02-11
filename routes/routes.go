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

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/user-signup", authController.Signup)
	http.HandleFunc("/resend-otp", authController.ResendOTP)
	http.HandleFunc("/verify-otp", authController.VerifyOTP)
	http.HandleFunc("/user-login", authController.Login)
	http.HandleFunc("/forgot-password", authController.ForgotPassword)
	http.HandleFunc("/forgot-otp", authController.ForgotOTP)
	http.HandleFunc("/reset-password", authController.ResetPassword)
	http.HandleFunc("/user-profile", middleware.AuthMiddleware(userController.GetProfile))
	http.HandleFunc("/update-profile", middleware.AuthMiddleware(userController.UpdateProfile))
	http.HandleFunc("/change-password", middleware.AuthMiddleware(userController.ChangePassword))
	http.HandleFunc("/create-book", middleware.AuthMiddleware(bookController.CreateBook))
	http.HandleFunc("/get-books", middleware.AuthMiddleware(bookController.GetBooks))
	http.HandleFunc("/update-book", middleware.AuthMiddleware(bookController.UpdateBook))
	http.HandleFunc("/delete-book", middleware.AuthMiddleware(bookController.DeleteBook))
	http.HandleFunc("/create-transaction", middleware.AuthMiddleware(transactionController.CreateTransaction))
	http.HandleFunc("/book-details", middleware.AuthMiddleware(transactionController.GetBookDetails))
	http.HandleFunc("/update-transaction", middleware.AuthMiddleware(transactionController.UpdateTransaction))
	http.HandleFunc("/delete-transaction", middleware.AuthMiddleware(transactionController.DeleteTransaction))
}
