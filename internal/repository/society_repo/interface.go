//go:generate mockgen -source=interface.go -destination=../../mock/mock_society_repo.go -package=mock
package society_repo

import "loopit/internal/models"

type SocietyRepo interface {
	Create(society models.Society) error
	FindAll() ([]models.Society, error)
	FindByID(id int) (models.Society, error)
	Save() error
}
