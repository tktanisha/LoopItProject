package handlers

import (
	"encoding/json"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/models"
	"loopit/internal/services/society_service"
	"loopit/pkg/logger"
	"net/http"
)

type SocietyHandler struct {
	societyService society_service.SocietyServiceInterface
	log            logger.LoggerInterface
}

func NewSocietyHandler(societyService society_service.SocietyServiceInterface, log logger.LoggerInterface) *SocietyHandler {
	return &SocietyHandler{societyService: societyService, log: log}
}

func (h *SocietyHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("GET /societies", h.GetAllSocieties)
	r.HandleFunc("POST /societies", h.CreateSociety)

}

func (h *SocietyHandler) GetAllSocieties(w http.ResponseWriter, r *http.Request) {
	societies, err := h.societyService.GetAllSocieties()
	if err != nil {
		h.log.Error("Failed to fetch societies: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to fetch societies",
			"error":   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    true,
		"societies": societies,
	})
}

func (h *SocietyHandler) CreateSociety(w http.ResponseWriter, r *http.Request) {
	userCtxVal := r.Context().Value(constants.UserCtxKey)
	if userCtxVal == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "unauthorized",
			"error":   "user context missing",
		})
		return
	}
	_, ok := userCtxVal.(*models.UserContext)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "internal error",
			"error":   "invalid user context",
		})
		return
	}

	var payload struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Pincode  string `json:"pincode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	//TODO:  Check if user is admin
	err := h.societyService.CreateSociety(payload.Name, payload.Location, payload.Pincode)
	if err != nil {
		h.log.Error("Failed to create society: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to create society",
			"error":   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "society created successfully",
	})
}
