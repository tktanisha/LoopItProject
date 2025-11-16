package handlers

import "loopit/internal/api/router"

//go:generate mockgen -source=handler.go -destination=../../mock/mock_handler.go -package=mock

type Handler interface {
	RegisterRoutes(router router.Router)
}
