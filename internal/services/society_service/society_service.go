package society_service

import (
	"loopit/internal/models"
	"loopit/internal/repository/society_repo"
	"time"
)

type SocietyService struct {
	societyRepo society_repo.SocietyRepo
}

func NewSocietyService(repo society_repo.SocietyRepo) SocietyServiceInterface {
	return &SocietyService{societyRepo: repo}
}

func (s *SocietyService) GetAllSocieties() ([]models.Society, error) {
	return s.societyRepo.FindAll()
}

func (s *SocietyService) CreateSociety(name, location, pincode string) error {
	society := models.Society{
		Name:      name,
		Location:  location,
		Pincode:   pincode,
		CreatedAt: time.Now(),
	}
	return s.societyRepo.Create(society)
}
