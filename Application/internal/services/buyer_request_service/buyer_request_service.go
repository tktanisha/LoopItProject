package buyer_request_service

import (
	"errors"
	"fmt"
	"loopit/internal/enums"
	br_status "loopit/internal/enums/buyer_request_status"
	order_status "loopit/internal/enums/order_status"
	"loopit/internal/models"
	"loopit/internal/repository/buyer_request_repo"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"

	"time"
)

type BuyerRequestService struct {
	buyerRequestRepo buyer_request_repo.BuyerRequestRepo
	productRepo      product_repo.ProductRepo
	orderRepo        order_repo.OrderRepo
	categoryRepo     category_repo.CategoryRepo
}

func NewBuyerRequestService(
	buyerReqRepo buyer_request_repo.BuyerRequestRepo,
	productRepo product_repo.ProductRepo,
	orderRepo order_repo.OrderRepo,
	categoryRepo category_repo.CategoryRepo,
) BuyerRequestServiceInterface {
	return &BuyerRequestService{
		buyerRequestRepo: buyerReqRepo,
		productRepo:      productRepo,
		orderRepo:        orderRepo,
		categoryRepo:     categoryRepo,
	}
}

func (s *BuyerRequestService) CreateBuyerRequest(productID int64, userCtx *models.UserContext) error {

	// Step 1: Validate the product
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("product not found")
	}
	if !product.Product.IsAvailable {
		return errors.New("product not available")
	}

	// Step 2: Prevent lender from requesting their own product
	if product.Product.LenderID == userCtx.ID {
		return errors.New("lender cannot create a buying request for their own product")
	}

	
	prodIDPtr := &productID
	statuses := []string{br_status.Pending.String()}

	// Pass productID and statuses to the repository for efficient filtering
	existingRequests, err := s.buyerRequestRepo.GetAllBuyerRequests(prodIDPtr, statuses)
	if err != nil {
		return err
	}

	for _, req := range existingRequests {
		if req.RequestedBy == userCtx.ID {
			return errors.New("a pending or approved request already exists")
		}
	}

	// Step 4: Create the new buyer request
	newRequest := models.BuyingRequest{
		ProductID:   productID,
		RequestedBy: userCtx.ID,
		Status:      br_status.Pending,
		CreatedAt:   time.Now(),
	}

	if err := s.buyerRequestRepo.CreateBuyerRequest(newRequest); err != nil {
		return err
	}

	return nil
}

func (s *BuyerRequestService) UpdateBuyerRequestStatus(requestID int64, updatedStatus br_status.Status, userCtx *models.UserContext) error {

	if userCtx.Role != enums.RoleLender {
		return errors.New("unauthorized: only lenders can update request status")
	}

	if updatedStatus != br_status.Approved && updatedStatus != br_status.Rejected {
		return errors.New("invalid status: only 'approved' or 'rejected' allowed")
	}

	allRequests, err := s.buyerRequestRepo.GetAllBuyerRequests(nil, nil)
	if err != nil {
		return err
	}

	var req *models.BuyingRequest
	for i := range allRequests {
		if allRequests[i].ID == requestID {
			req = &allRequests[i]
			break
		}
	}
	if req == nil {
		return errors.New("buyer request not found")
	}

	if updatedStatus == br_status.Rejected {
		if err := s.buyerRequestRepo.UpdateStatusBuyerRequest(requestID, br_status.Rejected.String()); err != nil {
			return err
		}
		return nil
	}

	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	category, err := s.categoryRepo.FindByID(product.Category.ID)
	if err != nil {
		return errors.New("category not found")
	}

	newOrder := models.Order{
		ProductID:      req.ProductID,
		UserID:         req.RequestedBy,
		StartDate:      time.Now(),
		EndDate:        time.Time{},
		TotalAmount:    category.Price,
		SecurityAmount: category.Security,
		Status:         order_status.InUse,
		CreatedAt:      time.Now(),
	}

	if err := s.orderRepo.CreateOrder(newOrder); err != nil {
		return err
	}

	if err := s.buyerRequestRepo.UpdateStatusBuyerRequest(requestID, br_status.Approved.String()); err != nil {
		return err
	}

	return nil
}

func (s *BuyerRequestService) GetAllBuyerRequests(productID *int64, filterStatuses []string) ([]models.BuyingRequest, error) {

	// Pass filters directly to the repository
	requests, err := s.buyerRequestRepo.GetAllBuyerRequests(productID, filterStatuses)
	if err != nil {
		return nil, err
	}

	fmt.Println("request=", requests)
	return requests, nil
}
