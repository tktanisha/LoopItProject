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
)

func InitDBRepos() error {
	LenderRepo = lender_repo.NewLenderDBRepo(db.DB)
	UserRepo = user_repo.NewUserDBRepo(db.DB, LenderRepo)
	CategoryRepo = category_repo.NewCategoryDBRepo(db.DB)
	ProductRepo = product_repo.NewProductDBRepo(db.DB, CategoryRepo, UserRepo)
	BuyerRequestRepo = buyer_request_repo.NewBuyerRequestDBRepo(db.DB)
	OrderRepo = order_repo.NewOrderDBRepo(db.DB, ProductRepo)
	ReturnRequestRepo = return_request_repo.NewReturnRequestDBRepo(db.DB)
	FeedBackRepo = feedback_repo.NewFeedBackDBRepo(db.DB)
	SocietyRepo = society_repo.NewSocietyDBRepo(db.DB)

	initServices()
	return nil
}
