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

var InitDBRepos = func(logger logger.LoggerInterface, db db.DatabaseInterface) error {

	// func InitDBRepos(logger logger.LoggerInterface, db db.DatabaseInterface) error {
	LenderRepo = lender_repo.NewLenderDBRepo(db, logger)
	UserRepo = user_repo.NewUserDBRepo(db, LenderRepo, logger)
	CategoryRepo = category_repo.NewCategoryDBRepo(db, logger)
	ProductRepo = product_repo.NewProductDBRepo(db, CategoryRepo, UserRepo, logger)
	BuyerRequestRepo = buyer_request_repo.NewBuyerRequestDBRepo(db, logger)
	OrderRepo = order_repo.NewOrderDBRepo(db, ProductRepo, logger)
	ReturnRequestRepo = return_request_repo.NewReturnRequestDBRepo(db, logger)
	FeedBackRepo = feedback_repo.NewFeedBackDBRepo(db, logger)
	SocietyRepo = society_repo.NewSocietyDBRepo(db, logger)

	InitServiceFiles(logger)
	return nil
}
