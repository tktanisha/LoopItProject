package main

import (
	"fmt"
	"loopit/internal/api"
	"loopit/internal/api/handlers"
	"loopit/internal/api/router"
	"loopit/internal/db"
	"loopit/internal/initializer"
	"loopit/pkg/logger"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	log := logger.GetLogger()
	defer log.Close()

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	DB_URL := os.Getenv("DB_URL")
	err = db.ConnectDB(DB_URL)

	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}

	initializer.InitServices()

	mux := http.NewServeMux()
	r := router.NewMuxRouter(mux)

	// create handlers
	userHandler := handlers.NewUserHandler(initializer.UserService, log)

	// register all
	api.SetupRoutes(r, userHandler)

	log.Info("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
