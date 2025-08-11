package utils

import (
	"context"
	"fmt"
	"loopit/internal/config"
	"loopit/internal/constants"
	"loopit/internal/models"
)

// GetAuthenticatedUserFromContext - Extract authenticated user context from the context
func GetAuthenticatedUserFromContext(ctx context.Context) (*models.UserContext, bool) {
	userCtxRaw := ctx.Value(constants.UserCtxKey)
	userCtx, ok := userCtxRaw.(*models.UserContext)
	if !ok || userCtx == nil {
		fmt.Println(config.Red + "Unauthorized access. Please login first." + config.Reset)
		return nil, false
	}
	return userCtx, true
}
