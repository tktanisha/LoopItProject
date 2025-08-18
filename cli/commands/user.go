package commands

import (
	"fmt"
	"loopit/cli/initializer"
	"loopit/internal/enums"
	"loopit/internal/models"
)

// Implement BecomeLender command
func BecomeLender(userCtx *models.UserContext) {
	log.Info(fmt.Sprintf("CLI: User %d attempting to become a lender", userCtx.ID))
	err := initializer.UserService.BecomeLender(userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error becoming a lender for user %d: %v", userCtx.ID, err))
		fmt.Println("Error becoming a lender:", err)
		return
	}

	// Update user context to reflect the new role
	userCtx.Role = enums.RoleLender

	log.Info(fmt.Sprintf("CLI: User %d has become a lender", userCtx.ID))
	fmt.Println("User is now a lender")
}
