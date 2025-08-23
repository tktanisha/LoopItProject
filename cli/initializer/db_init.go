package initializer

import (
	"loopit/internal/db"
	"loopit/internal/repository/buyer_request_repo"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/feedback_repo"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/return_request_repo"
	"loopit/internal/repository/society_repo"
	"loopit/internal/repository/user_repo"
	"loopit/pkg/logger"
)

func InitDBRepos(logger *logger.Logger) error {
	LenderRepo = lender_repo.NewLenderDBRepo(db.DB, logger)
	UserRepo = user_repo.NewUserDBRepo(db.DB, LenderRepo, logger)
	CategoryRepo = category_repo.NewCategoryDBRepo(db.DB, logger)
	ProductRepo = product_repo.NewProductDBRepo(db.DB, CategoryRepo, UserRepo, logger)
	BuyerRequestRepo = buyer_request_repo.NewBuyerRequestDBRepo(db.DB, logger)
	OrderRepo = order_repo.NewOrderDBRepo(db.DB, ProductRepo, logger)
	ReturnRequestRepo = return_request_repo.NewReturnRequestDBRepo(db.DB, logger)
	FeedBackRepo = feedback_repo.NewFeedBackDBRepo(db.DB, logger)
	SocietyRepo = society_repo.NewSocietyDBRepo(db.DB, logger)

	initServices(logger)
	return nil
}
