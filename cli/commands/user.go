package commands

import (
	"fmt"
	"loopit/internal/enums"
	"loopit/internal/models"
)

// Implement BecomeLender command
func BecomeLender(userCtx *models.UserContext) {
	err := UserService.BecomeLender(userCtx)

	if err != nil {
		fmt.Println("Error becoming a lender:", err)
		return

	}

	// Update user context to reflect the new role
	userCtx.Role = enums.RoleLender

	fmt.Println("User is now a lender")

}
