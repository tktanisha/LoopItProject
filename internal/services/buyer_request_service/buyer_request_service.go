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
	"loopit/pkg/logger"
	"time"
)

type BuyerRequestService struct {
	buyerRequestRepo buyer_request_repo.BuyerRequestRepo
	productRepo      product_repo.ProductRepo
	orderRepo        order_repo.OrderRepo
	categoryRepo     category_repo.CategoryRepo
	log              *logger.Logger
}

func NewBuyerRequestService(
	buyerReqRepo buyer_request_repo.BuyerRequestRepo,
	productRepo product_repo.ProductRepo,
	orderRepo order_repo.OrderRepo,
	categoryRepo category_repo.CategoryRepo,
	log *logger.Logger,
) BuyerRequestServiceInterface {
	return &BuyerRequestService{
		buyerRequestRepo: buyerReqRepo,
		productRepo:      productRepo,
		orderRepo:        orderRepo,
		categoryRepo:     categoryRepo,
		log:              log,
	}
}

func (s *BuyerRequestService) CreateBuyerRequest(productID int, userCtx *models.UserContext) error {
	s.log.Info(fmt.Sprintf("CreateBuyerRequest called by user %d for product %d", userCtx.ID, productID))

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		s.log.Warning(fmt.Sprintf("Product %d not found", productID))
		return errors.New("product not found")
	}
	if !product.Product.IsAvailable {
		s.log.Warning(fmt.Sprintf("Product %d is not available", productID))
		return errors.New("product not available")
	}

	if product.Product.LenderID == userCtx.ID {
		s.log.Warning(fmt.Sprintf("User %d attempted to request their own product %d", userCtx.ID, productID))
		return errors.New("lender cannot create a buying request for their own product")
	}

	allRequests, err := s.buyerRequestRepo.GetAllBuyerRequests([]string{br_status.Pending.String(), br_status.Approved.String()})
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch existing requests for user %d, error: %v", userCtx.ID, err))
		return err
	}
	for _, req := range allRequests {
		if req.ProductID == productID && req.RequestedBy == userCtx.ID {
			s.log.Warning(fmt.Sprintf("Duplicate buyer request found for user %d on product %d", userCtx.ID, productID))
			return errors.New("a pending or approved request already exists")
		}
	}

	newRequest := models.BuyingRequest{
		ProductID:   productID,
		RequestedBy: userCtx.ID,
		Status:      br_status.Pending,
		CreatedAt:   time.Now(),
	}

	if err := s.buyerRequestRepo.CreateBuyerRequest(newRequest); err != nil {
		s.log.Error(fmt.Sprintf("Failed to create buyer request for user %d, product %d, error: %v", userCtx.ID, productID, err))
		return err
	}

	s.log.Info(fmt.Sprintf("Buyer request created successfully for user %d, product %d", userCtx.ID, productID))
	return nil
}

func (s *BuyerRequestService) UpdateBuyerRequestStatus(requestID int, updatedStatus br_status.Status, userCtx *models.UserContext) error {
	s.log.Info(fmt.Sprintf("UpdateBuyerRequestStatus called by user %d for request %d to status %s", userCtx.ID, requestID, updatedStatus.String()))

	if userCtx.Role != enums.RoleLender {
		s.log.Warning(fmt.Sprintf("Unauthorized status update attempt by user %d", userCtx.ID))
		return errors.New("unauthorized: only lenders can update request status")
	}

	if updatedStatus != br_status.Approved && updatedStatus != br_status.Rejected {
		s.log.Warning(fmt.Sprintf("Invalid status %s attempted for request %d", updatedStatus.String(), requestID))
		return errors.New("invalid status: only 'approved' or 'rejected' allowed")
	}

	allRequests, err := s.buyerRequestRepo.GetAllBuyerRequests(nil)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch requests for status update by user %d, error: %v", userCtx.ID, err))
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
		s.log.Warning(fmt.Sprintf("Buyer request %d not found", requestID))
		return errors.New("buyer request not found")
	}

	if updatedStatus == br_status.Rejected {
		if err := s.buyerRequestRepo.UpdateStatusBuyerRequest(requestID, br_status.Rejected.String()); err != nil {
			s.log.Error(fmt.Sprintf("Failed to reject buyer request %d, error: %v", requestID, err))
			return err
		}
		s.log.Info(fmt.Sprintf("Buyer request %d rejected", requestID))
		return nil
	}

	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Product %d not found while approving request %d", req.ProductID, requestID))
		return errors.New("product not found")
	}

	category, err := s.categoryRepo.FindByID(product.Category.ID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Category %d not found for product %d, request %d", product.Category.ID, req.ProductID, requestID))
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
		s.log.Error(fmt.Sprintf("Failed to create order for request %d, error: %v", requestID, err))
		return err
	}

	if err := s.buyerRequestRepo.UpdateStatusBuyerRequest(requestID, br_status.Approved.String()); err != nil {
		s.log.Error(fmt.Sprintf("Failed to update buyer request %d to approved, error: %v", requestID, err))
		return err
	}

	s.log.Info(fmt.Sprintf("Buyer request %d approved and order created", requestID))
	return nil
}

func (s *BuyerRequestService) GetAllBuyerRequestsByStatus(productID int, status br_status.Status) ([]models.BuyingRequest, error) {
	s.log.Info(fmt.Sprintf("Fetching buyer requests for product %d with status %s", productID, status.String()))

	filtered, err := s.buyerRequestRepo.GetAllBuyerRequests([]string{status.String()})
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch buyer requests for product %d, error: %v", productID, err))
		return nil, err
	}

	result := []models.BuyingRequest{}
	for _, req := range filtered {
		if req.ProductID == productID {
			result = append(result, req)
		}
	}

	s.log.Info(fmt.Sprintf("Found %d buyer requests for product %d with status %s", len(result), productID, status.String()))
	return result, nil
}
