package commands

import (
	"fmt"
	"loopit/cli/initializer"
	"loopit/cli/utils"
	"loopit/internal/models"
	"os"

	"github.com/olekukonko/tablewriter"
)

func GiveFeedback(userCtx *models.UserContext) {
	orderId := utils.IntConversion(utils.Input("Enter your order ID"))
	feedbackText := utils.Input("Enter your feedback text")
	rating := utils.IntConversion(utils.Input("Enter your rating (1-5)"))

	log.Info(fmt.Sprintf("CLI: User %d attempting to give feedback for order %d", userCtx.ID, orderId))
	err := initializer.FeedBackService.GiveFeedback(orderId, feedbackText, rating, userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error giving feedback by user %d: %v", userCtx.ID, err))
		fmt.Printf("Error giving feedback: %v\n", err)
		return
	}

	log.Info(fmt.Sprintf("CLI: Feedback given successfully by user %d for order %d", userCtx.ID, orderId))
	fmt.Println("Feedback given successfully!")
}

func GetAllGivenFeedbacks(userCtx *models.UserContext) {
	log.Info(fmt.Sprintf("CLI: Fetching all feedbacks given by user %d", userCtx.ID))
	feedbacks, err := initializer.FeedBackService.GetAllGivenFeedbacks(userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching given feedbacks for user %d: %v", userCtx.ID, err))
		fmt.Println("Error fetching given feedbacks:", err)
		return
	}

	if len(feedbacks) == 0 {
		log.Info(fmt.Sprintf("CLI: No given feedbacks found for user %d", userCtx.ID))
		fmt.Println("No feedbacks given yet.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Given To", "Text", "Rating", "Created At")

	for _, f := range feedbacks {
		table.Append([]string{
			fmt.Sprintf("%d", f.ID),
			fmt.Sprintf("%d", f.GivenTo),
			f.Text,
			fmt.Sprintf("%d", f.Rating),
			f.CreatedAt.Format("2006-01-02"),
		})
	}
	table.Bulk(true)
	table.Render()
}

func GetAllReceivedFeedbacks(userCtx *models.UserContext) {
	log.Info(fmt.Sprintf("CLI: Fetching all feedbacks received by user %d", userCtx.ID))
	feedbacks, err := initializer.FeedBackService.GetAllReceivedFeedbacks(userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching received feedbacks for user %d: %v", userCtx.ID, err))
		fmt.Println("Error fetching received feedbacks:", err)
		return
	}
	if len(feedbacks) == 0 {
		log.Info(fmt.Sprintf("CLI: No received feedbacks found for user %d", userCtx.ID))
		fmt.Println("No feedbacks received yet.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Given By", "Text", "Rating", "Created At")
	for _, f := range feedbacks {
		table.Append([]string{
			fmt.Sprintf("%d", f.ID),
			fmt.Sprintf("%d", f.GivenBy),
			f.Text,
			fmt.Sprintf("%d", f.Rating),
			f.CreatedAt.Format("2006-01-02"),
		})
	}
	table.Bulk(true)
	table.Render()
}
