package product_repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/user_repo"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductDBRepo struct {
	db           *db.DynamoClient
	categoryRepo category_repo.CategoryRepo
	userRepo     user_repo.UserRepo
}

func NewProductDBRepo(db *db.DynamoClient, categoryRepo category_repo.CategoryRepo, userRepo user_repo.UserRepo) *ProductDBRepo {
	return &ProductDBRepo{
		db:           db,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}

func (r *ProductDBRepo) Create(product *models.Product) error {
	product.ID = time.Now().UnixNano()
	product.CreatedAt = time.Now()

	base := map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.ID)},
		"LenderID":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.LenderID)},
		"CategoryID":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.CategoryID)},
		"Name":        &types.AttributeValueMemberS{Value: product.Name},
		"Description": &types.AttributeValueMemberS{Value: product.Description},
		"Duration":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.Duration)},
		"IsAvailable": &types.AttributeValueMemberBOOL{Value: product.IsAvailable},
		"ImageUrl":    &types.AttributeValueMemberS{Value: product.ImageUrl},
		"CreatedAt":   &types.AttributeValueMemberS{Value: product.CreatedAt.Format(time.RFC3339)},
	}

	items := []map[string]types.AttributeValue{
		mergeMap(base, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)},
		}),
		mergeMap(base, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d#ID#%d", product.LenderID, product.ID)},
		}),
		mergeMap(base, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(product.Name), product.ID)},
		}),
		mergeMap(base, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("CATEGORY#%d", product.CategoryID)},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)},
		}),
	}

	for _, item := range items {
		_, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(r.db.Table),
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf("failed to create product item: %w", err)
		}
	}

	return nil
}

func (r *ProductDBRepo) FindByID(id int64) (*models.ProductResponse, error) {
	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", id)},
	}

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	if out.Item == nil {
		return nil, errors.New("product not found")
	}

	type productHelper struct {
		ID          int64     `dynamodbav:"ID"`
		LenderID    string    `dynamodbav:"LenderID"`
		CategoryID  string    `dynamodbav:"CategoryID"`
		Name        string    `dynamodbav:"Name"`
		Description string    `dynamodbav:"Description"`
		Duration    int       `dynamodbav:"Duration"`
		IsAvailable bool      `dynamodbav:"IsAvailable"`
		CreatedAt   time.Time `dynamodbav:"CreatedAt"`
		ImageUrl    string    `dynamodbav:"ImageUrl"`
	}

	var helper productHelper
	if err := attributevalue.UnmarshalMap(out.Item, &helper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	lenderID, _ := strconv.ParseInt(helper.LenderID, 10, 64)
	categoryID, _ := strconv.ParseInt(helper.CategoryID, 10, 64)

	var category models.Category
	if r.categoryRepo != nil {
		category, _ = r.categoryRepo.FindByID(categoryID)
	}

	var user models.User
	if r.userRepo != nil {
		u, err := r.userRepo.FindByID(lenderID)
		if err != nil {
			log.Printf("could not find user for lenderID=%d: %v", lenderID, err)
		}
		if u != nil {
			user = *u
		}
	}

	response := &models.ProductResponse{
		Product: models.Product{
			ID:          helper.ID,
			LenderID:    lenderID,
			CategoryID:  categoryID,
			Name:        helper.Name,
			Description: helper.Description,
			Duration:    helper.Duration,
			IsAvailable: helper.IsAvailable,
			CreatedAt:   helper.CreatedAt,
			ImageUrl:    helper.ImageUrl,
		},
		Category: category,
		User:     user,
	}

	log.Printf("FindByID result: %+v", response)
	return response, nil
}

func (r *ProductDBRepo) FindAll(filters models.ProductFilter) ([]*models.ProductResponse, error) {
	var pk, skPrefix string

	if filters.CategoryID != "" {
		pk = fmt.Sprintf("CATEGORY#%s", filters.CategoryID)
		skPrefix = "PRODUCT#"
	} else {
		pk = "PRODUCT"
		if filters.LenderID != "" {
			skPrefix = fmt.Sprintf("LENDER#%s", filters.LenderID)
		} else if filters.Search != "" {
			skPrefix = fmt.Sprintf("NAME#%s", strings.ToLower(filters.Search)) // normalized for prefix match
		} else {
			skPrefix = "PRODUCT#"
		}
	}

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":skPrefix": &types.AttributeValueMemberS{Value: skPrefix},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}

	type productHelper struct {
		ID          int64     `dynamodbav:"ID"`
		LenderID    string    `dynamodbav:"LenderID"`
		CategoryID  string    `dynamodbav:"CategoryID"`
		Name        string    `dynamodbav:"Name"`
		Description string    `dynamodbav:"Description"`
		Duration    int       `dynamodbav:"Duration"`
		IsAvailable bool      `dynamodbav:"IsAvailable"`
		CreatedAt   time.Time `dynamodbav:"CreatedAt"`
		ImageUrl    string    `dynamodbav:"ImageUrl"`
	}

	var helpers []productHelper
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &helpers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal products: %w", err)
	}

	var responses []*models.ProductResponse
	for _, h := range helpers {
		lenderID, _ := strconv.ParseInt(h.LenderID, 10, 64)
		categoryID, _ := strconv.ParseInt(h.CategoryID, 10, 64)

		var category models.Category
		if r.categoryRepo != nil {
			category, _ = r.categoryRepo.FindByID(categoryID)
		}

		var user models.User
		if r.userRepo != nil {
			u, err := r.userRepo.FindByID(lenderID)
			if err == nil && u != nil {
				user = *u
			}
		}

		responses = append(responses, &models.ProductResponse{
			Product: models.Product{
				ID:          h.ID,
				LenderID:    lenderID,
				CategoryID:  categoryID,
				Name:        h.Name,
				Description: h.Description,
				Duration:    h.Duration,
				IsAvailable: h.IsAvailable,
				CreatedAt:   h.CreatedAt,
				ImageUrl:    h.ImageUrl,
			},
			Category: category,
			User:     user,
		})
	}

	if filters.Search != "" {
		var filtered []*models.ProductResponse
		for _, p := range responses {
			if strings.Contains(strings.ToLower(p.Product.Name), strings.ToLower(filters.Search)) ||
				strings.Contains(strings.ToLower(p.Product.Description), strings.ToLower(filters.Search)) {
				filtered = append(filtered, p)
			}
		}
		return filtered, nil
	}

	if filters.IsAvailable != "" {
		var filtered []*models.ProductResponse
		wantAvailable := strings.ToLower(filters.IsAvailable) == "true"
		for _, p := range responses {
			if p.Product.IsAvailable == wantAvailable {
				filtered = append(filtered, p)
			}
		}
		return filtered, nil
	}

	return responses, nil
}

func (r *ProductDBRepo) Update(product *models.Product) error {
	ctx := context.TODO()

	existing, err := r.FindByID(product.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing product: %w", err)
	}

	if existing.Product.CategoryID != product.CategoryID {
		_, err := r.db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(r.db.Table),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("CATEGORY#%d", existing.Product.CategoryID)},
				"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete old category index: %w", err)
		}
	}

	if !strings.EqualFold(existing.Product.Name, product.Name) {
		_, err := r.db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(r.db.Table),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
				"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(existing.Product.Name), product.ID)},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete old name index: %w", err)
		}
	}

	updateExpr := "SET #n = :name, Description = :description, #dur = :duration, IsAvailable = :isAvailable, CategoryID = :categoryId, LenderID = :lenderId"
	exprAttrNames := map[string]string{"#n": "Name", "#dur": "Duration"}
	exprAttrValues := map[string]types.AttributeValue{
		":name":        &types.AttributeValueMemberS{Value: product.Name},
		":description": &types.AttributeValueMemberS{Value: product.Description},
		":duration":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.Duration)},
		":isAvailable": &types.AttributeValueMemberBOOL{Value: product.IsAvailable},
		":categoryId":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.CategoryID)},
		":lenderId":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.LenderID)},
	}

	keys := []map[string]types.AttributeValue{
		{"pk": &types.AttributeValueMemberS{Value: "PRODUCT"}, "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)}},
		{"pk": &types.AttributeValueMemberS{Value: "PRODUCT"}, "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d#ID#%d", product.LenderID, product.ID)}},
		{"pk": &types.AttributeValueMemberS{Value: "PRODUCT"}, "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(product.Name), product.ID)}},
		{"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("CATEGORY#%d", product.CategoryID)}, "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)}},
	}

	for _, key := range keys {
		_, err := r.db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(r.db.Table),
			Key:                       key,
			UpdateExpression:          aws.String(updateExpr),
			ExpressionAttributeNames:  exprAttrNames,
			ExpressionAttributeValues: exprAttrValues,
		})
		if err != nil {
			return fmt.Errorf("failed to update product record for key %+v: %w", key, err)
		}
	}

	return nil
}

func (r *ProductDBRepo) Delete(productID int64) error {
	ctx := context.TODO()

	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", productID)},
	}

	out, err := r.db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if out.Item == nil {
		return fmt.Errorf("product not found")
	}

	type productHelper struct {
		Name       string `dynamodbav:"Name"`
		LenderID   int64  `dynamodbav:"LenderID"`
		CategoryID int64  `dynamodbav:"CategoryID"`
	}
	var p productHelper
	if err := attributevalue.UnmarshalMap(out.Item, &p); err != nil {
		return fmt.Errorf("failed to unmarshal product: %w", err)
	}

	keys := []map[string]types.AttributeValue{
		{
			"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", productID)},
		},
		{
			"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d#ID#%d", p.LenderID, productID)},
		},
		{
			"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(p.Name), productID)},
		},
		{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("CATEGORY#%d", p.CategoryID)},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", productID)},
		},
	}

	for _, k := range keys {
		_, err := r.db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(r.db.Table),
			Key:       k,
		})
		if err != nil {
			return fmt.Errorf("failed to delete product record for key %+v: %w", k, err)
		}
	}

	log.Printf("Product deleted successfully from all indexes: ID=%d", productID)
	return nil
}

func (r *ProductDBRepo) Save() error {
	return nil
}

func mergeMap(base, extra map[string]types.AttributeValue) map[string]types.AttributeValue {
	merged := make(map[string]types.AttributeValue)
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range extra {
		merged[k] = v
	}
	return merged
}
