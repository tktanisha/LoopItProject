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
		fmt.Println(orderStatus, status)
		if err != nil {
			fmt.Println(config.Red+"Invalid status filter:"+config.Reset, err)
			return
		}
		filterStatus = append(filterStatus, orderStatus)
	}

	orders, err := initializer.OrderService.GetOrderHistory(userCtx, filterStatus)
	if err != nil {
		fmt.Println(config.Red+"Error fetching order history:"+config.Reset, err)
		return
	}

	if len(orders) == 0 {
		fmt.Println(config.Yellow + "No orders found." + config.Reset)
		return
	}

	printOrderTable("Order History", orders)
}

// 2. Mark Order as Returned (Only Lender)
func MarkOrderAsReturned(userCtx *models.UserContext) {
	if userCtx.Role != enums.RoleLender {
		fmt.Println(config.Red + "Only lenders can mark orders as returned." + config.Reset)
		return
	}

	orderID := utils.IntConversion(utils.Input("Enter Order ID to mark as returned: "))

	err := initializer.OrderService.MarkOrderAsReturned(orderID, userCtx)
	if err != nil {
		fmt.Println(config.Red+"Error marking order as returned:"+config.Reset, err)
		return
	}

	fmt.Println(config.Green + "Order marked as returned successfully!" + config.Reset)
}

// 3. Get All Approved Awaiting Orders (Only Lender)
func GetAllApprovedAwaitingOrders(userCtx *models.UserContext) {
	if userCtx.Role != enums.RoleLender {
		fmt.Println(config.Red + "Only lenders can view approved awaiting orders." + config.Reset)
		return
	}

	orders, err := initializer.OrderService.GetAllApprovedAwaitingOrders(userCtx)
	if err != nil {
		fmt.Println(config.Red+"Error fetching approved awaiting orders:"+config.Reset, err)
		return
	}

	if len(orders) == 0 {
		fmt.Println(config.Yellow + "No awaiting approved orders found." + config.Reset)
		return
	}

	printOrderTable("Approved Orders Awaiting Return", orders)
}

// 4. Get Lender Orders (Only Lender)
func GetLenderOrders(userCtx *models.UserContext) {
	if userCtx.Role != enums.RoleLender {
		fmt.Println(config.Red + "Only lenders can view their orders." + config.Reset)
		return
	}

	orders, err := initializer.OrderService.GetLenderOrders(userCtx)
	if err != nil {
		fmt.Println(config.Red+"Error fetching lender orders:"+config.Reset, err)
		return
	}

	if len(orders) == 0 {
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
