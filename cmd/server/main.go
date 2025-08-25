package main

import (
	"fmt"
	"loopit/internal/api"
	"loopit/internal/api/handlers"
	"loopit/internal/api/middleware"
	"loopit/internal/api/router"
	"loopit/internal/db"
	"loopit/internal/initializer"
	"loopit/pkg/logger"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Starting server...")
	log := logger.GetLogger()
	defer log.Close()

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	DB_URL := os.Getenv("DB_URL")
	err = db.ConnectDB(DB_URL)
	fmt.Println("Connecting to database...")

	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}

	initializer.InitServices()

	mux := http.NewServeMux()
	r := router.NewMuxRouter(mux)

	// create handlers
	userHandler := handlers.NewUserHandler(initializer.UserService, log)
	societyHandler := handlers.NewSocietyHandler(initializer.SocietyService, log)
	returnRequestHandler := handlers.NewReturnRequestHandler(initializer.ReturnRequestService, log)
	productHandler := handlers.NewProductHandler(initializer.ProductService, log)
	categoryHandler := handlers.NewCategoryHandler(initializer.CategoryService, log)
	authHandler := handlers.NewAuthHandler(initializer.AuthService, log)
	buyerRequestHandler := handlers.NewBuyerRequestHandler(initializer.BuyerRequestService, log)
	feedbackHandler := handlers.NewFeedbackHandler(initializer.FeedBackService, log)
	orderHandler := handlers.NewOrderHandler(initializer.OrderService, log)

	// register all
	api.SetupRoutes(r, userHandler, societyHandler, returnRequestHandler, productHandler, categoryHandler, authHandler, buyerRequestHandler, feedbackHandler, orderHandler)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		mux.ServeHTTP(w, req)
	})
	finalHandler := middleware.AuthMiddleware(log, protectedMux)

	log.Info("Server running on :8080")
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", finalHandler)
}
