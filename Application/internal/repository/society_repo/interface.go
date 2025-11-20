package society_repo

import "loopit/internal/models"

type SocietyRepo interface {
	Create(society models.Society) error
	FindAll() ([]models.Society, error)
	FindByID(id int64) (models.Society, error)
	Update(society models.Society) error
	Delete(id int64) error

}
