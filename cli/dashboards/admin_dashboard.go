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

// AdminDashboard - Main dashboard for administrators
func AdminDashboard(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintAdminMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			AdminSystemMenu(ctx, userCtx)
		case "2":
			fmt.Println(config.Yellow + "User Management - Coming Soon!" + config.Reset)
		case "3":
			commands.AuthLogout(&ctx)
			return
		case "4":
			fmt.Println("Exiting. Goodbye!")
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}

// AdminSystemMenu - System management functions for administrators
func AdminSystemMenu(ctx context.Context, userCtx *models.UserContext) {
	for {
		menus.PrintAdminSystemMenu()
		fmt.Print(config.Yellow + "Choose an option: " + config.Reset)
		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			commands.CreateSociety(userCtx)
		case "2":
			commands.CreateCategory()
		case "3":
			return
		default:
			fmt.Println(config.Red + "Invalid option. Try again." + config.Reset)
		}
	}
}
