package commands

import (
	"fmt"
	"loopit/cli/utils"
	"loopit/internal/models"
	"os"

	"github.com/olekukonko/tablewriter"
)

func GiveFeedback(userCtx *models.UserContext) {
	orderId := utils.IntConversion(utils.Input("Enter your order ID"))
	feedbackText := utils.Input("Enter your feedback text")
	rating := utils.IntConversion(utils.Input("Enter your rating (1-5)"))

	err := FeedBackService.GiveFeedback(orderId, feedbackText, rating, userCtx)
	if err != nil {
		fmt.Printf("Error giving feedback: %v\n", err)
		return
	}

	fmt.Println("Feedback given successfully!")
}

func GetAllGivenFeedbacks(userCtx *models.UserContext) {
	feedbacks, err := FeedBackService.GetAllGivenFeedbacks(userCtx)
	if err != nil {
		fmt.Println("Error fetching given feedbacks:", err)
		return
	}

	if len(feedbacks) == 0 {
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
	feedbacks, err := FeedBackService.GetAllReceivedFeedbacks(userCtx)
	if err != nil {
		fmt.Println("Error fetching received feedbacks:", err)
		return
	}
	if len(feedbacks) == 0 {
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
