package commands

import (
	"fmt"
	"loopit/cli/initializer"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/enums"
	"loopit/internal/enums/order_status"
	"loopit/internal/models"
	"os"

	"github.com/olekukonko/tablewriter"
)

// 1. Get Order History (Buyer or Lender)
func GetOrderHistory(userCtx *models.UserContext) {
	status := utils.Input("Enter status filter (In Use/Return Requested/Returned or leave empty): ")

	var filterStatus []order_status.Status
	if status != "" {
		orderStatus, err := order_status.ParseStatus(status)
		if err != nil {
			log.Warning(fmt.Sprintf("CLI: Invalid status filter input: %s by user %d", status, userCtx.ID))
			fmt.Println(config.Red+"Invalid status filter:"+config.Reset, err)
			return
		}
		filterStatus = append(filterStatus, orderStatus)
	}

	log.Info(fmt.Sprintf("CLI: User %d requests order history with status %v", userCtx.ID, filterStatus))
	orders, err := initializer.OrderService.GetOrderHistory(userCtx, filterStatus)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching order history for user %d: %v", userCtx.ID, err))
		fmt.Println(config.Red+"Error fetching order history:"+config.Reset, err)
		return
	}

	if len(orders) == 0 {
		log.Info(fmt.Sprintf("CLI: User %d found no orders matching filter", userCtx.ID))
		fmt.Println(config.Yellow + "No orders found." + config.Reset)
		return
	}

	printOrderTable("Order History", orders)
}

// 2. Mark Order as Returned (Only Lender)
func MarkOrderAsReturned(userCtx *models.UserContext) {
	if userCtx.Role != enums.RoleLender {
		log.Warning(fmt.Sprintf("CLI: Unauthorized mark as returned attempt by user %d (role: %s)", userCtx.ID, userCtx.Role))
		fmt.Println(config.Red + "Only lenders can mark orders as returned." + config.Reset)
		return
	}

	orderID := utils.IntConversion(utils.Input("Enter Order ID to mark as returned: "))
	log.Info(fmt.Sprintf("CLI: Lender %d marking order %d as returned", userCtx.ID, orderID))
	err := initializer.OrderService.MarkOrderAsReturned(orderID, userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error marking order %d as returned by user %d: %v", orderID, userCtx.ID, err))
		fmt.Println(config.Red+"Error marking order as returned:"+config.Reset, err)
		return
	}

	log.Info(fmt.Sprintf("CLI: Order %d marked as returned successfully by user %d", orderID, userCtx.ID))
	fmt.Println(config.Green + "Order marked as returned successfully!" + config.Reset)
}

// 3. Get All Approved Awaiting Orders (Only Lender)
func GetAllApprovedAwaitingOrders(userCtx *models.UserContext) {
	if userCtx.Role != enums.RoleLender {
		log.Warning(fmt.Sprintf("CLI: Unauthorized approved-awaiting order fetch attempt by user %d (role: %s)", userCtx.ID, userCtx.Role))
		fmt.Println(config.Red + "Only lenders can view approved awaiting orders." + config.Reset)
		return
	}

	log.Info(fmt.Sprintf("CLI: Lender %d fetching approved awaiting orders", userCtx.ID))
	orders, err := initializer.OrderService.GetAllApprovedAwaitingOrders(userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching approved awaiting orders for lender %d: %v", userCtx.ID, err))
		fmt.Println(config.Red+"Error fetching approved awaiting orders:"+config.Reset, err)
		return
	}

	if len(orders) == 0 {
		log.Info(fmt.Sprintf("CLI: Lender %d found no approved awaiting orders", userCtx.ID))
		fmt.Println(config.Yellow + "No awaiting approved orders found." + config.Reset)
		return
	}

	printOrderTable("Approved Orders Awaiting Return", orders)
}

// 4. Get Lender Orders (Only Lender)
func GetLenderOrders(userCtx *models.UserContext) {
	if userCtx.Role != enums.RoleLender {
		log.Warning(fmt.Sprintf("CLI: Unauthorized lender orders fetch attempt by user %d (role: %s)", userCtx.ID, userCtx.Role))
		fmt.Println(config.Red + "Only lenders can view their orders." + config.Reset)
		return
	}

	log.Info(fmt.Sprintf("CLI: Lender %d fetching their orders", userCtx.ID))
	orders, err := initializer.OrderService.GetLenderOrders(userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching lender orders for user %d: %v", userCtx.ID, err))
		fmt.Println(config.Red+"Error fetching lender orders:"+config.Reset, err)
		return
	}

	if len(orders) == 0 {
		log.Info(fmt.Sprintf("CLI: Lender %d found no orders", userCtx.ID))
		fmt.Println(config.Yellow + "No orders found for lender." + config.Reset)
		return
	}

	printOrderTable("Lender's Orders", orders)
}

func printOrderTable(title string, orders []*models.Order) {
	fmt.Println("\n" + title + ":")
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Order ID", "Product ID", "User ID", "Amount", "Security", "Status", "Start Date")

	for _, o := range orders {
		table.Append([]string{
			fmt.Sprintf("%d", o.ID),
			fmt.Sprintf("%d", o.ProductID),
			fmt.Sprintf("%d", o.UserID),
			fmt.Sprintf("%.2f", o.TotalAmount),
			fmt.Sprintf("%.2f", o.SecurityAmount),
			o.Status.String(),
			o.StartDate.Format("2006-01-02"),
		})
	}

	table.Bulk(true)
	table.Render()
}
