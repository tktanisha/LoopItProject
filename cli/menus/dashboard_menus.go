package menus

import (
	"fmt"
	"loopit/internal/config"
	"loopit/internal/enums"
)

// PrintUserMenu - Display the main user dashboard menu
func PrintUserMenu(role enums.Role) {
	fmt.Println(config.Green + "\n🛍️  USER DASHBOARD" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 🔍 Browse & Shop")
	fmt.Println("[2] 📦 My Orders & Requests")
	fmt.Println("[3] ⭐ Feedback & Reviews")
	fmt.Println("[4] 👤 Account Management")
	fmt.Println("[5] 🚪 Logout")
	fmt.Println("[6] ❌ Exit")

	if role == enums.RoleLender || role == enums.RoleAdmin {
		fmt.Println("[7] 🏪 Switch to Lender Dashboard")
	}

	if role == enums.RoleAdmin {
		fmt.Println("[8] ⚙️  Switch to Admin Dashboard")
	}
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintLenderMenu - Display the main lender dashboard menu
func PrintLenderMenu() {
	fmt.Println(config.Green + "\n🏪 LENDER DASHBOARD" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 📦 Product Management")
	fmt.Println("[2] 🛒 Order Management")
	fmt.Println("[3] 📋 Buyer Requests Management")
	fmt.Println("[4] ⭐ Feedback & Returns")
	fmt.Println("[5] ⬅️  Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintAdminMenu - Display the main admin dashboard menu
func PrintAdminMenu() {
	fmt.Println(config.Green + "\n⚙️  ADMIN DASHBOARD" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 🏗️  System Management")
	fmt.Println("[2] 👥 User Management")
	fmt.Println("[3] ⬅️  Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}
