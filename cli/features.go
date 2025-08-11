package cli

import (
	"context"
	"fmt"
	"loopit/cli/dashboards"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/enums"
)

// FeatureMenu - Main router that directs users to appropriate dashboards based on their role
func FeatureMenu(ctx context.Context) {
	userCtx, ok := utils.GetAuthenticatedUserFromContext(ctx)
	if !ok || userCtx == nil {
		fmt.Println(config.Red + "Unauthorized access. Please login first." + config.Reset)
		return
	}

	switch userCtx.Role {
	case enums.RoleUser:
		dashboards.UserDashboard(ctx, userCtx)

		userCtx, ok = utils.GetAuthenticatedUserFromContext(ctx)
		if ok && userCtx != nil && userCtx.Role != enums.RoleUser {
			FeatureMenu(ctx)
		}
	case enums.RoleLender:
		dashboards.LenderDashboard(ctx, userCtx)
	case enums.RoleAdmin:
		dashboards.AdminDashboard(ctx, userCtx)
	default:
		fmt.Println(config.Red + "Unknown user role." + config.Reset)
	}
}
