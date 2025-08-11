package buyer_request_repo

import (
	"errors"
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
	"loopit/internal/storage"
)

type BuyerRequestFileRepo struct {
	filePath       string
	buyingRequests []models.BuyingRequest
}

func NewBuyerRequestFileRepo(filePath string) (*BuyerRequestFileRepo, error) {
	requests, err := storage.ReadJSONFile[models.BuyingRequest](filePath)
	if err != nil {
		return nil, err
	}
	return &BuyerRequestFileRepo{
		filePath:       filePath,
		buyingRequests: requests,
	}, nil
}

func (r *BuyerRequestFileRepo) GetAllBuyerRequests(filterStatuses []string) ([]models.BuyingRequest, error) {
	if len(filterStatuses) == 0 {
		return r.buyingRequests, nil
	}

	statusMap := make(map[string]bool)
	for _, status := range filterStatuses {
		statusMap[status] = true
	}

	filtered := []models.BuyingRequest{}
	for _, req := range r.buyingRequests {
		if statusMap[req.Status.String()] {
			filtered = append(filtered, req)
		}
	}

	return filtered, nil
}

func (r *BuyerRequestFileRepo) UpdateStatusBuyerRequest(id int, newStatus string) error {
	for i, req := range r.buyingRequests {
		if req.ID == id {
			r.buyingRequests[i].Status, _ = buyer_request_status.ParseStatus(newStatus)
			return nil
		}
	}
	return errors.New("buying request not found")
}

func (r *BuyerRequestFileRepo) CreateBuyerRequest(req models.BuyingRequest) error {
	req.ID = len(r.buyingRequests) + 1
	r.buyingRequests = append(r.buyingRequests, req)
	return nil
}

func (r *BuyerRequestFileRepo) GetBuyerRequestByID(id int) (*models.BuyingRequest, error) {
	for _, req := range r.buyingRequests {
		if req.ID == id {
			return &req, nil
		}
	}
	return nil, errors.New("buying request not found")
}

func (r *BuyerRequestFileRepo) Save() error {
	return storage.WriteJSONFile(r.filePath, r.buyingRequests)
}
