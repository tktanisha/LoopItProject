// package product_repo

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"loopit/internal/db"
// 	"loopit/internal/models"
// 	"loopit/internal/repository/category_repo"
// 	"loopit/internal/repository/user_repo"
// 	"loopit/pkg/logger"
// 	"time"
// )

// type ProductDBRepo struct {
// 	db           *db.DynamoClient
// 	categoryRepo category_repo.CategoryRepo
// 	userRepo     user_repo.UserRepo
// 	log          logger.LoggerInterface
// }

// func NewProductDBRepo(db *db.DynamoClient, categoryRepo category_repo.CategoryRepo, userRepo user_repo.UserRepo) *ProductDBRepo {
// 	return &ProductDBRepo{
// 		db:           db,
// 		categoryRepo: categoryRepo,
// 		userRepo:     userRepo,
// 	}
// }

// // FindAll returns all products with category and lender info
// func (r *ProductDBRepo) FindAll(filters models.ProductFilter) ([]*models.ProductResponse, error) {
// 	query := `
// 		SELECT id, lender_id, category_id, name, description, duration, is_available, created_at, image_url
// 		FROM products WHERE 1=1
// 	`
// 	args := []any{}
// 	argIdx := 1 // For PostgreSQL placeholders

// 	if filters.Search != "" {
// 		query += fmt.Sprintf(" AND (LOWER(name) LIKE LOWER($%d) OR LOWER(description) LIKE LOWER($%d))", argIdx, argIdx+1)
// 		searchParam := "%" + filters.Search + "%"
// 		args = append(args, searchParam, searchParam)
// 		argIdx += 2
// 	}
// 	if filters.LenderID != "" {
// 		query += fmt.Sprintf(" AND lender_id = $%d", argIdx)
// 		args = append(args, filters.LenderID)
// 		argIdx++
// 	}
// 	if filters.CategoryID != "" {
// 		query += fmt.Sprintf(" AND category_id = $%d", argIdx)
// 		args = append(args, filters.CategoryID)
// 		argIdx++
// 	}
// 	if filters.IsAvailable != "" {
// 		query += fmt.Sprintf(" AND is_available = $%d", argIdx)
// 		args = append(args, filters.IsAvailable)
// 		argIdx++
// 	}

// 	r.log.Info(fmt.Sprintf("Executing product query: %s with args %v", query, args))

// 	rows, err := r.db.Query(query, args...)
// 	if err != nil {
// 		if r.log != nil {
// 			r.log.Error(fmt.Sprintf("DB error fetching products: %v", err))
// 		}
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var responses []*models.ProductResponse
// 	for rows.Next() {
// 		var p models.Product
// 		if err := rows.Scan(&p.ID, &p.LenderID, &p.CategoryID, &p.Name, &p.Description, &p.Duration, &p.IsAvailable, &p.CreatedAt, &p.ImageUrl); err != nil {
// 			if r.log != nil {
// 				r.log.Warning(fmt.Sprintf("DB warning scanning product row: %v", err))
// 			}
// 			continue
// 		}

// 		category, _ := r.categoryRepo.FindByID(p.CategoryID)
// 		user, _ := r.userRepo.FindByID(p.LenderID)

// 		responses = append(responses, &models.ProductResponse{
// 			Product:  p,
// 			Category: category,
// 			User:     *user,
// 		})
// 	}

// 	return responses, nil
// }

// // FindByID returns a single product by ID with category and lender info
// func (r *ProductDBRepo) FindByID(id int64) (*models.ProductResponse, error) {
// 	row := r.db.QueryRow("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products WHERE id=$1", id)

// 	var p models.Product
// 	if err := row.Scan(&p.ID, &p.LenderID, &p.CategoryID, &p.Name, &p.Description, &p.Duration, &p.IsAvailable, &p.CreatedAt); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			if r.log != nil {
// 				r.log.Warning(fmt.Sprintf("DB: No product found with id %d", id))
// 			}
// 			return nil, errors.New("product not found")
// 		}
// 		if r.log != nil {
// 			r.log.Error(fmt.Sprintf("DB error fetching product by id %d: %v", id, err))
// 		}
// 		return nil, err
// 	}

// 	category, _ := r.categoryRepo.FindByID(p.CategoryID)
// 	user, _ := r.userRepo.FindByID(p.LenderID)

// 	return &models.ProductResponse{
// 		Product:  p,
// 		Category: category,
// 		User:     *user,
// 	}, nil
// }

// // Create inserts a new product into the database
// func (r *ProductDBRepo) Create(product *models.Product) error {
// 	query := `
//     INSERT INTO products (lender_id, category_id, name, description, duration, is_available, image_url, created_at)
//     VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
//     RETURNING id
//     `
// 	err := r.db.QueryRow(query, product.LenderID, product.CategoryID, product.Name, product.Description, product.Duration, product.IsAvailable, product.ImageUrl, time.Now()).Scan(&product.ID)
// 	if err != nil && r.log != nil {
// 		r.log.Error(fmt.Sprintf("DB error creating product '%s': %v", product.Name, err))
// 	}
// 	return err
// }

// // update
// func (r *ProductDBRepo) Update(product *models.Product) error {
// 	query := `
// 	UPDATE products
// 	SET lender_id=$1, category_id=$2, name=$3, description=$4, duration=$5, is_available=$6
// 	WHERE id=$7
// 	`
// 	_, err := r.db.Exec(query, product.LenderID, product.CategoryID, product.Name, product.Description, product.Duration, product.IsAvailable, product.ID)
// 	if err != nil && r.log != nil {
// 		r.log.Error(fmt.Sprintf("DB error updating product ID %d: %v", product.ID, err))
// 	}
// 	return err
// }

// // delete
// func (r *ProductDBRepo) Delete(id int64) error {
// 	_, err := r.db.Exec("DELETE FROM products WHERE id=$1", id)
// 	if err != nil && r.log != nil {
// 		r.log.Error(fmt.Sprintf("DB error deleting product ID %d: %v", id, err))
// 	}
// 	return err
// }

// // Save is a no-op for Postgres as changes are applied immediately
// func (r *ProductDBRepo) Save() error {
// 	return nil
// }

package product_repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/user_repo"

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

func NewProductDBRepo(
	db *db.DynamoClient,
	categoryRepo category_repo.CategoryRepo,
	userRepo user_repo.UserRepo,
) *ProductDBRepo {
	return &ProductDBRepo{
		db:           db,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}

func (r *ProductDBRepo) Create(p *models.Product) error {
	p.ID = time.Now().UnixNano()
	p.CreatedAt = time.Now()

	table := r.db.Table

	pkBase := "PRODUCT"
	skByID := fmt.Sprintf("PRODUCT#%d", p.ID)
	skByLender := fmt.Sprintf("LENDER#%d#ID#%d", p.LenderID, p.ID)
	skByName := fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(p.Name), p.ID)

	// Common Item
	item := map[string]interface{}{
		"ID":          p.ID,
		"LenderID":    p.LenderID,
		"CategoryID":  p.CategoryID,
		"Name":        p.Name,
		"Description": p.Description,
		"Duration":    p.Duration,
		"IsAvailable": p.IsAvailable,
		"CreatedAt":   p.CreatedAt.Format(time.RFC3339),
		"ImageURL":    p.ImageUrl,
	}

	// Write item for each access pattern
	skValues := []string{skByID, skByLender, skByName}

	for _, sk := range skValues {
		marshalled, err := attributevalue.MarshalMap(
			map[string]interface{}{
				"pk": pkBase,
				"sk": sk,
			},
		)
		if err != nil {
			return fmt.Errorf("marshal key failed: %w", err)
		}

		body, err := attributevalue.MarshalMap(item)
		if err != nil {
			return fmt.Errorf("marshal body failed: %w", err)
		}

		for k, v := range body {
			marshalled[k] = v
		}

		_, err = r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item:      marshalled,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ProductDBRepo) FindByID(id int64) (*models.ProductResponse, error) {
	pk := "PRODUCT"
	sk := fmt.Sprintf("PRODUCT#%d", id)

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})

	if err != nil || out.Item == nil {
		return nil, fmt.Errorf("product not found")
	}

	var product models.Product
	if err := attributevalue.UnmarshalMap(out.Item, &product); err != nil {
		return nil, err
	}

	category, _ := r.categoryRepo.FindByID(product.CategoryID)
	user, _ := r.userRepo.FindByID(product.LenderID)

	return &models.ProductResponse{
		Product:  product,
		Category: category,
		User:     *user,
	}, nil
}

func (r *ProductDBRepo) FindByLender(lenderID int64) ([]models.ProductResponse, error) {
	pk := "PRODUCT"
	prefix := fmt.Sprintf("LENDER#%d#ID#", lenderID)

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
			":prefix": &types.AttributeValueMemberS{Value: prefix},
		},
	})

	if err != nil {
		return nil, err
	}

	var products []models.Product

	if err := attributevalue.UnmarshalListOfMaps(out.Items, &products); err != nil {
		return nil, err
	}

	// Convert to product responses
	var responses []models.ProductResponse
	for _, p := range products {
		c, _ := r.categoryRepo.FindByID(p.CategoryID)
		u, _ := r.userRepo.FindByID(p.LenderID)

		responses = append(responses, models.ProductResponse{
			Product:  p,
			Category: c,
			User:     *u,
		})
	}

	return responses, nil
}

func (r *ProductDBRepo) FindByName(name string) ([]models.ProductResponse, error) {
	pk := "PRODUCT"
	prefix := fmt.Sprintf("NAME#%s", strings.ToLower(name))

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
			":prefix": &types.AttributeValueMemberS{Value: prefix},
		},
	})

	if err != nil {
		return nil, err
	}

	var products []models.Product
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &products); err != nil {
		return nil, err
	}

	var responses []models.ProductResponse
	for _, p := range products {
		c, _ := r.categoryRepo.FindByID(p.CategoryID)
		u, _ := r.userRepo.FindByID(p.LenderID)

		responses = append(responses, models.ProductResponse{
			Product:  p,
			Category: c,
			User:     *u,
		})
	}

	return responses, nil
}

func (r *ProductDBRepo) Delete(id int64, lenderID int64, name string) error {
	pk := "PRODUCT"

	sks := []string{
		fmt.Sprintf("PRODUCT#%d", id),
		fmt.Sprintf("LENDER#%d#ID#%d", lenderID, id),
		fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(name), id),
	}

	for _, sk := range sks {
		_, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String(r.db.Table),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: pk},
				"sk": &types.AttributeValueMemberS{Value: sk},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ProductDBRepo) Save() error {
	return nil
}
