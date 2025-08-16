package main

import (
	"context"
	"fmt"
	"loopit/cli"
	"loopit/cli/commands"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/db"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

func cleanup() {
	fmt.Println("Saving all the files...")
	commands.UserFileRepo.Save()
	commands.ProductFileRepo.Save()
	commands.LenderFileRepo.Save()
	commands.CategoryFileRepo.Save()
	commands.BuyerRequestFileRepo.Save()
	commands.OrderFileRepo.Save()
	commands.ReturnRequestFileRepo.Save()
	commands.FeedBackFileRepo.Save()
	commands.SocietyFileRepo.Save()
}

func main() {
	// Setup cleanup on interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nReceived interrupt signal. Cleaning up...")
		cleanup()
		os.Exit(1)
	}()
	defer cleanup()

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

	commands.InitServices()
	utils.ShowBanner()

	ctx := context.Background()

	for {
		fmt.Println()
		fmt.Println("[1] Login")
		fmt.Println("[2] Create Account")
		fmt.Println("[3] Exit")
		fmt.Println()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)

		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			if commands.AuthLogin(&ctx) {
				cli.FeatureMenu(ctx)
			}
		case "2":
			if commands.AuthRegister(&ctx) {
				cli.FeatureMenu(ctx)
			}
		case "3":
			fmt.Println("Exiting. Goodbye! ")
			return
		default:
			fmt.Println(config.Red + "Invalid choice. Try again." + config.Reset)
		}
	}
}
