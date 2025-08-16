package initializer

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
)

func InitFileRepos() error {
	var err error

	LenderRepo, err = lender_repo.NewLenderFileRepo("data/lenders.json")
	if err != nil {
		return err
	}

	UserRepo, err = user_repo.NewUserFileRepo("data/users.json", LenderRepo)
	if err != nil {
		return err
	}

	CategoryRepo, err = category_repo.NewCategoryFileRepo("data/category.json")
	if err != nil {
		return err
	}

	ProductRepo, err = product_repo.NewProductFileRepo("data/products.json", CategoryRepo, UserRepo)
	if err != nil {
		return err
	}

	BuyerRequestRepo, err = buyer_request_repo.NewBuyerRequestFileRepo("data/buyer_requests.json")
	if err != nil {
		return err
	}

	OrderRepo, err = order_repo.NewOrderFileRepo("data/orders.json", ProductRepo)
	if err != nil {
		return err
	}

	ReturnRequestRepo, err = return_request_repo.NewReturnRequestFileRepo("data/return_requests.json")
	if err != nil {
		return err
	}

	FeedBackRepo, err = feedback_repo.NewFeedBackFileRepo("data/feedbacks.json")
	if err != nil {
		return err
	}

	SocietyRepo, err = society_repo.NewSocietyFileRepo("data/societies.json")
	if err != nil {
		return err
	}

	initServices()
	return nil
}
