package dashboards

import (
	"context"
	"fmt"
	"loopit/cli/commands"
	"loopit/cli/menus"
	"loopit/internal/config"
	"loopit/internal/models"
	"strings"
)

// UserDashboard - Main dashboard for regular users
func UserDashboard(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintUserMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			BrowsingMenu(ctx, userCtx)
		case "2":
			UserOrdersMenu(ctx, userCtx)
		case "3":
			FeedbackMenu(ctx, userCtx)
		case "4":
			AccountMenu(ctx, userCtx)
		case "5":
			commands.AuthLogout(&ctx)
			return
		case "6":
			fmt.Println("Exiting. Goodbye!")
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// BrowsingMenu - Browse and shopping functionality
func BrowsingMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintBrowsingMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			commands.GetAllProducts()
		case "2":
			commands.GetProductByID()
		case "3":
			commands.GetAllCategories()
		case "4":
			commands.GetAllSocieties()
		case "5":
			commands.CreateBuyerRequest(userCtx)
		case "6":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// UserOrdersMenu - Order and request management for users
func UserOrdersMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintUserOrdersMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			commands.GetOrderHistory(userCtx)
		case "2":
			commands.UpdateReturnRequestStatus(userCtx)
		case "3":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// FeedbackMenu - Feedback and review management
func FeedbackMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintFeedbackMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			commands.GiveFeedback(userCtx)
		case "2":
			commands.GetAllGivenFeedbacks(userCtx)
		case "3":
			commands.GetAllReceivedFeedbacks(userCtx)
		case "4":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// AccountMenu - Account management for users
func AccountMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintAccountMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			commands.BecomeLender(userCtx)
			fmt.Println(config.Green + "\n🎉 Congratulations! You are now a lender!" + config.Reset)
			fmt.Println(config.Yellow + "Redirecting to Lender Dashboard..." + config.Reset)
			return
		case "2":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}
