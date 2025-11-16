package society_service

import (
	"log"
	"loopit/internal/models"
	"loopit/internal/repository/society_repo"
	"time"
)

type SocietyService struct {
	societyRepo society_repo.SocietyRepo

}

func NewSocietyService(repo society_repo.SocietyRepo) SocietyServiceInterface {
	return &SocietyService{
		societyRepo: repo,
	}
}

func (s *SocietyService) GetAllSocieties() ([]models.Society, error) {
	societies, err := s.societyRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return societies, nil
}

func (s *SocietyService) CreateSociety(name, location, pincode string) error {

	society := models.Society{
		Name:      name,
		Location:  location,
		Pincode:   pincode,
		CreatedAt: time.Now(),
	}

	err := s.societyRepo.Create(society)
	if err != nil {
		return err
	}

	return nil
}

func (s *SocietyService) UpdateSociety(id int64, name, location, pincode string) error {
	log.Print("entered in service")
	society, err := s.societyRepo.FindByID(id)
	if err != nil {

		return err
	}
	log.Print("after repo ")
	society.Name = name
	society.Location = location
	society.Pincode = pincode
	err = s.societyRepo.Update(society)
	if err != nil {
		return err
	}
	log.Print("after society update")
	return nil
}

func (s *SocietyService) DeleteSociety(id int64) error {
	err := s.societyRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

