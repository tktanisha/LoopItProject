package api

import (
	"loopit/internal/api/handlers"
	"loopit/internal/api/router"
)

func SetupRoutes(r router.Router, hs ...handlers.Handler) {
	for _, h := range hs {
		h.RegisterRoutes(r)
	}
}
