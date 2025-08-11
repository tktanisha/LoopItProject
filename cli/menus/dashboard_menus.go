package menus

import (
	"fmt"
	"loopit/internal/config"
)

// PrintUserMenu - Display the main user dashboard menu
func PrintUserMenu() {
	fmt.Println(config.Green + "\n🛍️  USER DASHBOARD" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 🔍 Browse & Shop")
	fmt.Println("[2] 📦 My Orders & Requests")
	fmt.Println("[3] ⭐ Feedback & Reviews")
	fmt.Println("[4] 👤 Account Management")
	fmt.Println("[5] 🚪 Logout")
	fmt.Println("[6] ❌ Exit")
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
	fmt.Println("[5] 🔍 Browse as Customer")
	fmt.Println("[6] 🚪 Logout")
	fmt.Println("[7] ❌ Exit")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintAdminMenu - Display the main admin dashboard menu
func PrintAdminMenu() {
	fmt.Println(config.Green + "\n⚙️  ADMIN DASHBOARD" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 🏗️  System Management")
	fmt.Println("[2] 👥 User Management")
	fmt.Println("[3] 🏪 Browse as Lender")
	fmt.Println("[4] 🚪 Logout")
	fmt.Println("[5] ❌ Exit")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}
