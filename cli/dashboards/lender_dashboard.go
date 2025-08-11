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

// LenderDashboard - Main dashboard for lenders
func LenderDashboard(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintLenderMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			// Product Management
			LenderProductMenu(ctx, userCtx)
		case "2":
			// Order Management
			LenderOrderMenu(ctx, userCtx)
		case "3":
			// Buyer Requests Management
			LenderBuyerRequestMenu(ctx, userCtx)
		case "4":
			// Feedback & Returns
			LenderFeedbackMenu(ctx, userCtx)
		case "5":
			// Browse as customer
			BrowsingMenu(ctx, userCtx)
		case "6":
			commands.AuthLogout(&ctx)
			return
		case "7":
			fmt.Println("Exiting. Goodbye!")
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// LenderProductMenu - Product management for lenders
func LenderProductMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintLenderProductMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			commands.CreateProduct(userCtx)
			fmt.Println(config.Green + "\n✅ Product created successfully!" + config.Reset)
			fmt.Println(config.Yellow + "What would you like to do next?" + config.Reset)
			ProductActionMenu(ctx, userCtx)
		case "2":
			commands.GetAllProducts()
		case "3":
			commands.GetProductByID()
		case "4":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// LenderOrderMenu - Order management for lenders
func LenderOrderMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintLenderOrderMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			commands.GetLenderOrders(userCtx)
		case "2":
			commands.GetAllApprovedAwaitingOrders(userCtx)
		case "3":
			commands.CreateReturnRequest(userCtx)
		case "4":
			commands.MarkOrderAsReturned(userCtx)
		case "5":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// LenderBuyerRequestMenu - Buyer request management for lenders
func LenderBuyerRequestMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintLenderBuyerRequestMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			commands.GetAllBuyerRequests()
		case "2":
			commands.UpdateBuyerRequestStatus(userCtx)
		case "3":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// LenderFeedbackMenu - Feedback and returns for lenders
func LenderFeedbackMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintLenderFeedbackMenu()
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
			commands.CreateReturnRequest(userCtx)
		case "5":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// ProductActionMenu - Context-aware menu after creating a product
func ProductActionMenu(ctx context.Context, userCtx *models.UserContext) {
	fmt.Println(config.Green + "\n🚀 NEXT ACTIONS" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] ➕ Create Another Product")
	fmt.Println("[2] 📋 Manage Buyer Requests")
	fmt.Println("[3] 🛒 View My Orders")
	fmt.Println("[4] 📱 View All My Products")
	fmt.Println("[5] ⬅️  Back to Product Management")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)

	fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
	var choice string
	fmt.Scanln(&choice)

	switch strings.TrimSpace(choice) {
	case "1":
		commands.CreateProduct(userCtx)
		ProductActionMenu(ctx, userCtx)
	case "2":
		LenderBuyerRequestMenu(ctx, userCtx)
	case "3":
		LenderOrderMenu(ctx, userCtx)
	case "4":
		commands.GetAllProducts()
	case "5":
		return
	default:
		fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		ProductActionMenu(ctx, userCtx)
	}
}
