package menus

import (
	"fmt"
	"loopit/internal/config"
)

// User Submenus

// PrintBrowsingMenu - Display the browsing and shopping menu
func PrintBrowsingMenu() {
	fmt.Println(config.Green + "\n🔍 BROWSE & SHOP" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 📱 View All Products")
	fmt.Println("[2] 🔎 Search Product by ID")
	fmt.Println("[3] 📋 View All Categories")
	fmt.Println("[4] 🏘️ View All Societies")
	fmt.Println("[5] 📝 Create Buyer Request")
	fmt.Println("[6] ⬅️ Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintUserOrdersMenu - Display the user orders and requests menu
func PrintUserOrdersMenu() {
	fmt.Println(config.Green + "\n📦 ORDERS & REQUESTS" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 📜 View Order History")
	fmt.Println("[2] 🔄 Update Return Request Status")
	fmt.Println("[3] ⬅️ Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintFeedbackMenu - Display the feedback and reviews menu
func PrintFeedbackMenu() {
	fmt.Println(config.Green + "\n⭐ FEEDBACK" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 💬 Give Feedback")
	fmt.Println("[2] 📤 Given Feedbacks")
	fmt.Println("[3] 📥 Received Feedbacks")
	fmt.Println("[4] ⬅️ Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintAccountMenu - Display the account management menu
func PrintAccountMenu() {
	fmt.Println(config.Green + "\n👤 ACCOUNT MANAGEMENT" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 🏪  Become a Lender")
	fmt.Println("[2] ⬅️  Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// Lender Submenus

// PrintLenderProductMenu - Display the lender product management menu
func PrintLenderProductMenu() {
	fmt.Println(config.Green + "\n📦 PRODUCT MANAGEMENT" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] ➕ Create New Product")
	fmt.Println("[2] 📱 View All Products")
	fmt.Println("[3] 🔎 Search Product by ID")
	fmt.Println("[4] ⬅️  Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintLenderOrderMenu - Display the lender order management menu
func PrintLenderOrderMenu() {
	fmt.Println(config.Green + "\n🛒 ORDER MANAGEMENT" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 📜 View Lending History")
	fmt.Println("[2] ⏳ View Approved Awaiting Orders")
	fmt.Println("[3] ✅ Create Return Request")
	fmt.Println("[4] ✅ Mark Order as Returned")
	fmt.Println("[5] ⬅️ Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintLenderBuyerRequestMenu - Display the lender buyer request management menu
func PrintLenderBuyerRequestMenu() {
	fmt.Println(config.Green + "\n📋 BUYER REQUESTS MANAGEMENT" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 📄  View All Buyer Requests")
	fmt.Println("[2] ✏️  Update Buyer Request Status")
	fmt.Println("[3] ⬅️  Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// PrintLenderFeedbackMenu - Display the lender feedback and returns menu
func PrintLenderFeedbackMenu() {
	fmt.Println(config.Green + "\n⭐ FEEDBACK " + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 💬  Give Feedback")
	fmt.Println("[2] 📤  Given Feedbacks")
	fmt.Println("[3] 📥  Received Feedbacks")
	fmt.Println("[4] ⬅️  Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}

// Admin Submenus

// PrintAdminSystemMenu - Display the admin system management menu
func PrintAdminSystemMenu() {
	fmt.Println(config.Green + "\n🏗️  SYSTEM MANAGEMENT" + config.Reset)
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
	fmt.Println("[1] 🏘️  Create Society")
	fmt.Println("[2] 📂  Create Category")
	fmt.Println("[3] ⬅️  Back")
	fmt.Println(config.Cyan + "═══════════════════════════════════════" + config.Reset)
}
