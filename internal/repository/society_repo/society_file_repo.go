package society_repo

import (
	"errors"
	"loopit/internal/models"
	"loopit/internal/storage"
)

type SocietyFileRepo struct {
	societyFile string
	societies   []models.Society
}

func NewSocietyFileRepo(societyFile string) (*SocietyFileRepo, error) {
	societies, err := storage.ReadJSONFile[models.Society](societyFile)
	if err != nil {
		return nil, err
	}
	return &SocietyFileRepo{
		societyFile: societyFile,
		societies:   societies,
	}, nil
}

func (r *SocietyFileRepo) Create(society models.Society) error {
	id := len(r.societies)
	society.ID = id + 1
	r.societies = append(r.societies, society)
	return nil
}

func (r *SocietyFileRepo) FindAll() ([]models.Society, error) {
	return r.societies, nil
}

func (r *SocietyFileRepo) FindByID(id int) (models.Society, error) {
	for _, s := range r.societies {
		if s.ID == id {
			return s, nil
		}
	}
	return models.Society{}, errors.New("society not found")
}

func (r *SocietyFileRepo) Save() error {
	return storage.WriteJSONFile(r.societyFile, r.societies)
}
