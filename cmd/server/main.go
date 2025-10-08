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
	"github.com/rs/cors"
)

func main() {

	fmt.Println("Starting server...")
	var log logger.LoggerInterface = logger.GetLogger()
	defer log.Close()

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	DB_URL := os.Getenv("DB_URL")
	pg, err := db.ConnectDB(DB_URL)
	fmt.Println("Connecting to database...")

	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}

	fmt.Println("hello")
	// err = db.ExecuteSQLFile(pg, "internal/db/init_db_table.sql")
	// if err != nil {
	// 	log.Fatal(fmt.Sprintf("Error initializing tables: %v", err))
	// }

	initializer.InitServices(pg, log)
	fmt.Println("hi")
	// create handlers
	userHandler := handlers.NewUserHandler(initializer.UserService, log)
	societyHandler := handlers.NewSocietyHandler(initializer.SocietyService, log)
	returnRequestHandler := handlers.NewReturnRequestHandler(initializer.ReturnRequestService, log)
	productHandler := handlers.NewProductHandler(initializer.ProductService, log)
	categoryHandler := handlers.NewCategoryHandler(initializer.CategoryService, log)
	authHandler := handlers.NewAuthHandler(initializer.AuthService, log)
	buyerRequestHandler := handlers.NewBuyerRequestHandler(initializer.BuyerRequestService, initializer.ProductService, log)
	feedbackHandler := handlers.NewFeedbackHandler(initializer.FeedBackService, log)
	orderHandler := handlers.NewOrderHandler(initializer.OrderService, initializer.ProductService, log)

	publicMux := http.NewServeMux()
	publicRouter := router.NewMuxRouter(publicMux)
	protectedMux := http.NewServeMux()
	protectedRouter := router.NewMuxRouter(protectedMux)

	api.SetupRoutes(publicRouter, authHandler)
	api.SetupRoutes(protectedRouter, userHandler, societyHandler, returnRequestHandler, productHandler, categoryHandler, buyerRequestHandler, feedbackHandler, orderHandler)

	finalHandler := http.NewServeMux()
	finalHandler.Handle("/auth/", publicMux)
	finalHandler.Handle("/", middleware.AuthMiddleware(log, protectedMux))
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})
	handler := c.Handler(finalHandler)

	log.Info("Server running on :8080")
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", handler)
}
