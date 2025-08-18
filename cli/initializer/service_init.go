package initializer

import (
	"loopit/internal/services/auth_service"
	"loopit/internal/services/buyer_request_service"
	"loopit/internal/services/category_service"
	"loopit/internal/services/feedback_service"
	"loopit/internal/services/order_service"
	"loopit/internal/services/product_service"
	"loopit/internal/services/return_request_service"
	"loopit/internal/services/society_service"
	"loopit/internal/services/user_service"
	"loopit/pkg/logger"
)

func initServices(logger *logger.Logger) {
	AuthService = auth_service.NewAuthService(UserRepo, logger)
	ProductService = product_service.NewProductService(ProductRepo, UserRepo, logger)
	CategoryService = category_service.NewCategoryService(CategoryRepo, logger)
	UserService = user_service.NewUserService(UserRepo, logger)
	BuyerRequestService = buyer_request_service.NewBuyerRequestService(BuyerRequestRepo, ProductRepo, OrderRepo, CategoryRepo, logger)
	OrderService = order_service.NewOrderService(OrderRepo, ReturnRequestRepo, ProductRepo, logger)
	ReturnRequestService = return_request_service.NewReturnRequestService(OrderRepo, ProductRepo, ReturnRequestRepo, logger)
	FeedBackService = feedback_service.NewFeedbackService(FeedBackRepo, ProductRepo, OrderRepo, logger)
	SocietyService = society_service.NewSocietyService(SocietyRepo, logger)
}
