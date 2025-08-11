package commands

import (
	"fmt"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

// 1. Create Return Request (by lender)
func CreateReturnRequest(userCtx *models.UserContext) {
	orderID := utils.IntConversion(utils.Input("Enter Order ID to return: "))

	err := ReturnRequestService.CreateReturnRequest(userCtx.ID, orderID)
	if err != nil {
		fmt.Println(config.Red+"Error creating return request:"+config.Reset, err)
		return
	}

	fmt.Println(config.Green + "Return request created successfully!" + config.Reset)
}

// 2. Get all pending Return Requests (for user)
func GetAllPendingReturnRequests(userCtx *models.UserContext) {

	requests, err := ReturnRequestService.GetPendingReturnRequests(userCtx.ID)
	if err != nil {
		fmt.Println(config.Red+"Error fetching pending return requests:"+config.Reset, err)
		return
	}

	fmt.Println("\nPending Return Requests:")
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Order ID", "Requested By", "Status", "Created At")

	for _, r := range requests {
		table.Append([]string{
			fmt.Sprintf("%d", r.ID),
			fmt.Sprintf("%d", r.OrderID),
			fmt.Sprintf("%d", r.RequestedBy),
			r.Status.String(),
			r.CreatedAt.Format(time.RFC822),
		})
	}
	table.Render()
}

// 3. Update Return Request Status (accept/reject by user who placed the order)
func UpdateReturnRequestStatus(userCtx *models.UserContext) {
	reqID := utils.IntConversion(utils.Input("Enter Return Request ID to update: "))

	statusOptions := []string{return_request_status.Approved.String(), return_request_status.Rejected.String()}
	_, selectedStatus := utils.SelectFromList("Select new status", statusOptions)
	if selectedStatus == "" {
		fmt.Println(config.Red + "Status selection cancelled." + config.Reset)
		return
	}

	newStatus, err := return_request_status.ParseStatus(selectedStatus)
	if err != nil {
		fmt.Println(config.Red+"Error parsing status:"+config.Reset, err)
		return
	}

	err = ReturnRequestService.UpdateReturnRequestStatus(userCtx.ID, reqID, newStatus)
	if err != nil {
		fmt.Println(config.Red+"Error updating return request status:"+config.Reset, err)
		return
	}

	fmt.Println(config.Green + "Return request status updated to '" + selectedStatus + "' successfully!" + config.Reset)
}
