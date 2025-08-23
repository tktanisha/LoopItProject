package society_service

import (
	"fmt"
	"loopit/internal/models"
	"loopit/internal/repository/society_repo"
	"loopit/pkg/logger"
	"time"
)

type SocietyService struct {
	societyRepo society_repo.SocietyRepo
	log         *logger.Logger
}

func NewSocietyService(repo society_repo.SocietyRepo, log *logger.Logger) SocietyServiceInterface {
	return &SocietyService{
		societyRepo: repo,
		log:         log,
	}
}

func (s *SocietyService) GetAllSocieties() ([]models.Society, error) {
	s.log.Info("Service: Fetching all societies")
	societies, err := s.societyRepo.FindAll()
	if err != nil {
		s.log.Error("Service: Failed to fetch societies: " + err.Error())
		return nil, err
	}
	s.log.Info(fmt.Sprintf("Service: Successfully fetched %d societies", len(societies)))
	return societies, nil
}

func (s *SocietyService) CreateSociety(name, location, pincode string) error {
	s.log.Info(fmt.Sprintf("Service: Creating new society (name=%s)", name))

	society := models.Society{
		Name:      name,
		Location:  location,
		Pincode:   pincode,
		CreatedAt: time.Now(),
	}

	err := s.societyRepo.Create(society)
	if err != nil {
		s.log.Error(fmt.Sprintf("Service: Failed to create society (name=%s): %v", name, err))
		return err
	}

	s.log.Info(fmt.Sprintf("Service: Society created successfully (name=%s)", society.Name))
	return nil
}
