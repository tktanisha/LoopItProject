package cli

import (
	"context"
	"fmt"
	"loopit/cli/dashboards"
	"loopit/cli/utils"
	"loopit/internal/config"
)

// FeatureMenu - Main router that directs users to appropriate dashboards based on their role
func FeatureMenu(ctx context.Context) {
	userCtx, ok := utils.GetAuthenticatedUserFromContext(ctx)
	if !ok || userCtx == nil {
		fmt.Println(config.Red + "Unauthorized access. Please login first." + config.Reset)
		return
	}
	dashboards.UserDashboard(ctx, userCtx)
}
