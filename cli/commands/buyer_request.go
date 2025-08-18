package commands

import (
	"fmt"
	"loopit/cli/initializer"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
	"os"

	"github.com/olekukonko/tablewriter"
)

// 1. Buy a Product (Create a Buyer Request)
func CreateBuyerRequest(userCtx *models.UserContext) {
	productID := utils.IntConversion(utils.Input("Enter Product ID to buy: "))
	log.Info(fmt.Sprintf("CLI: CreateBuyerRequest invoked by user %d for product %d", userCtx.ID, productID))

	err := initializer.BuyerRequestService.CreateBuyerRequest(productID, userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Failed to create buyer request (user_id=%d, product_id=%d): %v", userCtx.ID, productID, err))
		fmt.Println(config.Red+"Error creating buyer request:"+config.Reset, err)
		return
	}

	log.Info(fmt.Sprintf("CLI: Buyer request created successfully (user_id=%d, product_id=%d)", userCtx.ID, productID))
	fmt.Println(config.Green + "Buyer request created successfully!" + config.Reset)
}

// 2. Get all buyer requests (status = pending, approved, rejected)
func GetAllBuyerRequests() {
	productId := utils.IntConversion(utils.Input("Enter Product ID to fetch buyer requests: "))
	log.Info(fmt.Sprintf("CLI: GetAllBuyerRequests invoked for product_id=%d", productId))

	requests, err := initializer.BuyerRequestService.GetAllBuyerRequestsByStatus(productId, buyer_request_status.Pending)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Failed to fetch buyer requests for product_id=%d: %v", productId, err))
		fmt.Println(config.Red+"Error fetching buyer requests:"+config.Reset, err)
		return
	}

	if len(requests) == 0 {
		log.Info(fmt.Sprintf("CLI: No buyer requests found for product_id=%d", productId))
		fmt.Println(config.Yellow + "No buyer requests found." + config.Reset)
		return
	}

	fmt.Println("\nBuyer Requests:")

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Product ID", "Requested By", "Status", "Created At")

	for _, r := range requests {
		table.Append([]string{
			fmt.Sprintf("%d", r.ID),
			fmt.Sprintf("%d", r.ProductID),
			fmt.Sprintf("%d", r.RequestedBy),
			r.Status.String(),
			r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	table.Bulk(true)
	table.Render()
	log.Info(fmt.Sprintf("CLI: Displayed %d buyer requests for product_id=%d", len(requests), productId))
}

// 3. Update Buyer Request Status
func UpdateBuyerRequestStatus(userCtx *models.UserContext) {
	reqID := utils.IntConversion(utils.Input("Enter Buyer Request ID to update: "))
	log.Info(fmt.Sprintf("CLI: UpdateBuyerRequestStatus invoked by user %d for request_id=%d", userCtx.ID, reqID))

	statusOptions := []string{buyer_request_status.Approved.String(), buyer_request_status.Rejected.String()}
	_, selectedStatus := utils.SelectFromList("Select new status", statusOptions)
	if selectedStatus == "" {
		log.Warning(fmt.Sprintf("CLI: Status selection cancelled for request_id=%d", reqID))
		fmt.Println(config.Red + "Status selection cancelled." + config.Reset)
		return
	}

	selectedStatusEnum, err := buyer_request_status.ParseStatus(selectedStatus)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Invalid status '%s' selected for request_id=%d: %v", selectedStatus, reqID, err))
		fmt.Println(config.Red+"Error parsing status:"+config.Reset, err)
		return
	}

	err = initializer.BuyerRequestService.UpdateBuyerRequestStatus(reqID, selectedStatusEnum, userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Failed to update buyer request %d to status '%s': %v", reqID, selectedStatus, err))
		fmt.Println(config.Red+"Error updating status:"+config.Reset, err)
		return
	}

	log.Info(fmt.Sprintf("CLI: Buyer request %d updated to status '%s' successfully", reqID, selectedStatus))
	fmt.Println(config.Green + "Buyer request status updated to '" + selectedStatus + "' successfully!" + config.Reset)
}
