package handlers

import "loopit/internal/api/router"

type Handler interface {
	RegisterRoutes(router router.Router)
}
