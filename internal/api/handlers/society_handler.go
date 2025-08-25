package handlers

import (
	"loopit/internal/api/router"
	"loopit/internal/services/society_service"
	"loopit/pkg/logger"
	"net/http"
)

type SocietyHandler struct {
	societyService society_service.SocietyServiceInterface
	log            *logger.Logger
}

func NewSocietyHandler(societyService society_service.SocietyServiceInterface, log *logger.Logger) *SocietyHandler {
	return &SocietyHandler{societyService: societyService, log: log}
}

func (h *SocietyHandler) RegisterRoutes(r router.Router) {
	r.Handle("/societies", http.HandlerFunc(h.GetAllSocieties))
}

func (h *SocietyHandler) GetAllSocieties(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Fetching all societies")
}
