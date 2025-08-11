package commands

import (
	"loopit/internal/repository/buyer_request_repo"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/feedback_repo"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/return_request_repo"
	"loopit/internal/repository/society_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/auth_service"
	"loopit/internal/services/buyer_request_service"
	"loopit/internal/services/category_service"
	"loopit/internal/services/feedback_service"
	"loopit/internal/services/order_service"
	"loopit/internal/services/product_service"
	"loopit/internal/services/return_request_service"
	"loopit/internal/services/society_service"
	"loopit/internal/services/user_service"
)

var (
	LenderFileRepo        lender_repo.LenderRepo
	UserFileRepo          user_repo.UserRepo
	ProductFileRepo       product_repo.ProductRepo
	CategoryFileRepo      category_repo.CategoryRepo
	BuyerRequestFileRepo  buyer_request_repo.BuyerRequestRepo
	OrderFileRepo         order_repo.OrderRepo
	ReturnRequestFileRepo return_request_repo.ReturnRequestRepo
	FeedBackFileRepo      feedback_repo.FeedbackRepository
	SocietyFileRepo       society_repo.SocietyRepo

	AuthService          auth_service.AuthServiceInterface
	ProductService       product_service.ProductServiceInterface
	CategoryService      category_service.CategoryServiceInterface
	UserService          user_service.UserServiceInterface
	BuyerRequestService  buyer_request_service.BuyerRequestServiceInterface
	OrderService         order_service.OrderServiceInterface
	ReturnRequestService return_request_service.ReturnRequestServiceInterface
	FeedBackService      feedback_service.FeedbackServiceInterface
	SocietyService       society_service.SocietyServiceInterface
)

func InitServices() error {
	var err error

	LenderFileRepo, err = lender_repo.NewLenderFileRepo("data/lenders.json")
	if err != nil {
		return err
	}

	UserFileRepo, err = user_repo.NewUserFileRepo("data/users.json", LenderFileRepo)
	if err != nil {
		return err
	}

	CategoryFileRepo, err = category_repo.NewCategoryFileRepo("data/category.json")
	if err != nil {
		return err
	}

	ProductFileRepo, err = product_repo.NewProductFileRepo("data/products.json", CategoryFileRepo, UserFileRepo)
	if err != nil {
		return err
	}

	BuyerRequestFileRepo, err = buyer_request_repo.NewBuyerRequestFileRepo("data/buyer_requests.json")
	if err != nil {
		return err
	}

	OrderFileRepo, err = order_repo.NewOrderFileRepo("data/orders.json", ProductFileRepo)
	if err != nil {
		return err
	}

	ReturnRequestFileRepo, err = return_request_repo.NewReturnRequestFileRepo("data/return_requests.json")
	if err != nil {
		return err
	}

	FeedBackFileRepo, err = feedback_repo.NewFeedBackFileRepo("data/feedbacks.json")
	if err != nil {
		return err
	}

	SocietyFileRepo, err = society_repo.NewSocietyFileRepo("data/societies.json")
	if err != nil {
		return err
	}

	AuthService = auth_service.NewAuthService(UserFileRepo)
	ProductService = product_service.NewProductService(ProductFileRepo, UserFileRepo)
	CategoryService = category_service.NewCategoryService(CategoryFileRepo)
	UserService = user_service.NewUserService(UserFileRepo)
	BuyerRequestService = buyer_request_service.NewBuyerRequestService(BuyerRequestFileRepo, ProductFileRepo, OrderFileRepo, CategoryFileRepo)
	OrderService = order_service.NewOrderService(OrderFileRepo, ReturnRequestFileRepo, ProductFileRepo)
	ReturnRequestService = return_request_service.NewReturnRequestService(OrderFileRepo, ProductFileRepo, ReturnRequestFileRepo)
	FeedBackService = feedback_service.NewFeedbackService(FeedBackFileRepo, ProductFileRepo, OrderFileRepo)
	SocietyService = society_service.NewSocietyService(SocietyFileRepo)

	return nil
}
