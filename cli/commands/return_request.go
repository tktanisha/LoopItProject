package commands

import (
	"fmt"
	"loopit/cli/initializer"
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

	log.Info(fmt.Sprintf("CLI: User %d initiates return request for order %d", userCtx.ID, orderID))
	err := initializer.ReturnRequestService.CreateReturnRequest(userCtx.ID, orderID)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Failed to create return request: %v", err))
		fmt.Println(fmt.Sprintf("%sError creating return request:%s %v", config.Red, config.Reset, err))
		return
	}

	log.Info(fmt.Sprintf("CLI: Return request created for order %d by user %d", orderID, userCtx.ID))
	fmt.Println(fmt.Sprintf("%sReturn request created successfully!%s", config.Green, config.Reset))
}

// 2. Get all pending Return Requests (for user)
func GetAllPendingReturnRequests(userCtx *models.UserContext) {
	log.Info(fmt.Sprintf("CLI: Fetching pending return requests for user %d", userCtx.ID))
	requests, err := initializer.ReturnRequestService.GetPendingReturnRequests(userCtx.ID)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Failed to get pending return requests: %v", err))
		fmt.Println(fmt.Sprintf("%serror fetching pending return requests:%s %v", config.Red, config.Reset, err))
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
		fmt.Println(fmt.Sprintf("%sStatus selection cancelled.%s", config.Red, config.Reset))
		return
	}

	newStatus, err := return_request_status.ParseStatus(selectedStatus)
	if err != nil {
		fmt.Println(fmt.Sprintf("%sError parsing status:%s %v", config.Red, config.Reset, err))
		return
	}

	log.Info(fmt.Sprintf("CLI: User %d attempts to update return request %d to status %s", userCtx.ID, reqID, newStatus.String()))
	err = initializer.ReturnRequestService.UpdateReturnRequestStatus(userCtx.ID, reqID, newStatus)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Failed to update return request status: %v", err))
		fmt.Println(fmt.Sprintf("%sError updating return request status:%s %v", config.Red, config.Reset, err))
		return
	}

	log.Info(fmt.Sprintf("CLI: Return request %d updated to '%s' by user %d", reqID, selectedStatus, userCtx.ID))
	fmt.Println(fmt.Sprintf("%sReturn request status updated to '%s' successfully!%s", config.Green, selectedStatus, config.Reset))
}
