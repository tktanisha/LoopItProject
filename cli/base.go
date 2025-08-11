package cli

import (
	"context"
	"fmt"
	"loopit/cli/commands"
	"loopit/cli/utils"
	"loopit/internal/config"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

func StartCLI() {
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

	commands.InitServices()
	utils.ShowBanner()

	ctx := context.Background()

	// userCtx, ok := utils.GetAuthenticatedUserFromContext(ctx)
	// if ok && userCtx != nil {
	// 	FeatureMenu(ctx)
	// 	return
	// }

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
				FeatureMenu(ctx)
			}
		case "2":
			if commands.AuthRegister(&ctx) {
				FeatureMenu(ctx)
			}
		case "3":
			fmt.Println("Exiting. Goodbye! ")
			return
		default:
			fmt.Println(config.Red + "Invalid choice. Try again." + config.Reset)
		}
	}
}
