package initializer

import (
	"loopit/internal/config"
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
	// Repo instances
	LenderRepo        lender_repo.LenderRepo
	UserRepo          user_repo.UserRepo
	ProductRepo       product_repo.ProductRepo
	CategoryRepo      category_repo.CategoryRepo
	BuyerRequestRepo  buyer_request_repo.BuyerRequestRepo
	OrderRepo         order_repo.OrderRepo
	ReturnRequestRepo return_request_repo.ReturnRequestRepo
	FeedBackRepo      feedback_repo.FeedbackRepository
	SocietyRepo       society_repo.SocietyRepo

	// Services
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
	if config.AppConfig.StorageType == "db" {
		return InitDBRepos()
	}
	return InitFileRepos()
}
