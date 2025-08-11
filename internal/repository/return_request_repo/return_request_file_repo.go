package return_request_repo

import (
	"errors"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/internal/storage"
)

type ReturnRequestFileRepo struct {
	filePath       string
	returnRequests []models.ReturnRequest
}

func NewReturnRequestFileRepo(filePath string) (*ReturnRequestFileRepo, error) {
	requests, err := storage.ReadJSONFile[models.ReturnRequest](filePath)
	if err != nil {
		return nil, err
	}
	return &ReturnRequestFileRepo{
		filePath:       filePath,
		returnRequests: requests,
	}, nil
}

func (r *ReturnRequestFileRepo) CreateReturnRequest(req models.ReturnRequest) error {
	req.ID = len(r.returnRequests) + 1
	r.returnRequests = append(r.returnRequests, req)
	return nil
}

func (r *ReturnRequestFileRepo) UpdateReturnRequestStatus(id int, newStatus string) error {
	for i, rr := range r.returnRequests {
		if rr.ID == id {
			newStatusEnum, err := return_request_status.ParseStatus(newStatus)
			if err != nil {
				return err
			}
			r.returnRequests[i].Status = newStatusEnum
			return nil
		}
	}
	return errors.New("return request not found")
}

func (r *ReturnRequestFileRepo) GetAllReturnRequests(filterStatuses []string) ([]models.ReturnRequest, error) {
	return r.returnRequests, nil
	if len(filterStatuses) == 0 {
		return r.returnRequests, nil
	}

	statusMap := make(map[string]bool)
	for _, status := range filterStatuses {
		statusMap[status] = true
	}

	filtered := []models.ReturnRequest{}
	for _, rr := range r.returnRequests {
		if statusMap[rr.Status.String()] {
			filtered = append(filtered, rr)
		}
	}

	return filtered, nil
}

func (r *ReturnRequestFileRepo) GetReturnRequestByID(id int) (models.ReturnRequest, error) {
	for _, rr := range r.returnRequests {
		if rr.ID == id {
			return rr, nil
		}
	}
	return models.ReturnRequest{}, errors.New("return request not found")
}

func (r *ReturnRequestFileRepo) Save() error {
	return storage.WriteJSONFile(r.filePath, r.returnRequests)
}
