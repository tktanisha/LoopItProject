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
)

func initServices() {
	AuthService = auth_service.NewAuthService(UserRepo)
	ProductService = product_service.NewProductService(ProductRepo, UserRepo)
	CategoryService = category_service.NewCategoryService(CategoryRepo)
	UserService = user_service.NewUserService(UserRepo)
	BuyerRequestService = buyer_request_service.NewBuyerRequestService(BuyerRequestRepo, ProductRepo, OrderRepo, CategoryRepo)
	OrderService = order_service.NewOrderService(OrderRepo, ReturnRequestRepo, ProductRepo)
	ReturnRequestService = return_request_service.NewReturnRequestService(OrderRepo, ProductRepo, ReturnRequestRepo)
	FeedBackService = feedback_service.NewFeedbackService(FeedBackRepo, ProductRepo, OrderRepo)
	SocietyService = society_service.NewSocietyService(SocietyRepo)
}
