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
	log         logger.LoggerInterface
}

func NewSocietyService(repo society_repo.SocietyRepo, log logger.LoggerInterface) SocietyServiceInterface {
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

func (s *SocietyService) UpdateSociety(id int, name, location, pincode string) error {
	s.log.Info(fmt.Sprintf("Service: Updating society (id=%d)", id))
	society, err := s.societyRepo.FindByID(id)
	if err != nil {
		s.log.Error(fmt.Sprintf("Service: Society not found (id=%d): %v", id, err))
		return err
	}
	society.Name = name
	society.Location = location
	society.Pincode = pincode
	err = s.societyRepo.Update(society)
	if err != nil {
		s.log.Error(fmt.Sprintf("Service: Failed to update society (id=%d): %v", id, err))
		return err
	}
	s.log.Info(fmt.Sprintf("Service: Society updated successfully (id=%d)", id))
	return nil
}

func (s *SocietyService) DeleteSociety(id int) error {
	s.log.Info(fmt.Sprintf("Service: Deleting society (id=%d)", id))
	err := s.societyRepo.Delete(id)
	if err != nil {
		s.log.Error(fmt.Sprintf("Service: Failed to delete society (id=%d): %v", id, err))
		return err
	}
	s.log.Info(fmt.Sprintf("Service: Society deleted successfully (id=%d)", id))
	return nil
}

