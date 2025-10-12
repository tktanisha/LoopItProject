package society_service

import "loopit/internal/models"

type SocietyServiceInterface interface {
	GetAllSocieties() ([]models.Society, error)
	CreateSociety(name, location, pincode string) error
	UpdateSociety(id int, name, location, pincode string) error
	DeleteSociety(id int) error
}
