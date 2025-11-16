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

// package product_repo

// import (
// 	"context"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"loopit/internal/db"
// 	"loopit/internal/models"
// 	"loopit/internal/repository/category_repo"
// 	"loopit/internal/repository/user_repo"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
// )

// type ProductDBRepo struct {
// 	db           *db.DynamoClient
// 	categoryRepo category_repo.CategoryRepo
// 	userRepo     user_repo.UserRepo
// }

// func NewProductDBRepo(
// 	db *db.DynamoClient,
// 	categoryRepo category_repo.CategoryRepo,
// 	userRepo user_repo.UserRepo,
// ) *ProductDBRepo {
// 	return &ProductDBRepo{
// 		db:           db,
// 		categoryRepo: categoryRepo,
// 		userRepo:     userRepo,
// 	}
// }

// func (r *ProductDBRepo) Create(p *models.Product) error {
// 	p.ID = time.Now().UnixNano()
// 	p.CreatedAt = time.Now()

// 	table := r.db.Table

// 	pkBase := "PRODUCT"
// 	skByID := fmt.Sprintf("PRODUCT#%d", p.ID)
// 	skByLender := fmt.Sprintf("LENDER#%d#ID#%d", p.LenderID, p.ID)
// 	skByName := fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(p.Name), p.ID)

// 	// Common Item
// 	item := map[string]interface{}{
// 		"ID":          p.ID,
// 		"LenderID":    p.LenderID,
// 		"CategoryID":  p.CategoryID,
// 		"Name":        p.Name,
// 		"Description": p.Description,
// 		"Duration":    p.Duration,
// 		"IsAvailable": p.IsAvailable,
// 		"CreatedAt":   p.CreatedAt.Format(time.RFC3339),
// 		"ImageURL":    p.ImageUrl,
// 	}

// 	// Write item for each access pattern
// 	skValues := []string{skByID, skByLender, skByName}

// 	for _, sk := range skValues {
// 		marshalled, err := attributevalue.MarshalMap(
// 			map[string]interface{}{
// 				"pk": pkBase,
// 				"sk": sk,
// 			},
// 		)
// 		if err != nil {
// 			return fmt.Errorf("marshal key failed: %w", err)
// 		}

// 		body, err := attributevalue.MarshalMap(item)
// 		if err != nil {
// 			return fmt.Errorf("marshal body failed: %w", err)
// 		}

// 		for k, v := range body {
// 			marshalled[k] = v
// 		}

// 		_, err = r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
// 			TableName: aws.String(table),
// 			Item:      marshalled,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (r *ProductDBRepo) FindByID(id int64) (*models.ProductResponse, error) {
// 	pk := "PRODUCT"
// 	sk := fmt.Sprintf("PRODUCT#%d", id)

// 	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
// 		TableName: aws.String(r.db.Table),
// 		Key: map[string]types.AttributeValue{
// 			"pk": &types.AttributeValueMemberS{Value: pk},
// 			"sk": &types.AttributeValueMemberS{Value: sk},
// 		},
// 	})

// 	if err != nil || out.Item == nil {
// 		return nil, fmt.Errorf("product not found")
// 	}

// 	var product models.Product
// 	if err := attributevalue.UnmarshalMap(out.Item, &product); err != nil {
// 		return nil, err
// 	}

// 	category, _ := r.categoryRepo.FindByID(product.CategoryID)
// 	user, _ := r.userRepo.FindByID(product.LenderID)

// 	return &models.ProductResponse{
// 		Product:  product,
// 		Category: category,
// 		User:     *user,
// 	}, nil
// }

// func (r *ProductDBRepo) FindByLender(lenderID int64) ([]models.ProductResponse, error) {
// 	pk := "PRODUCT"
// 	prefix := fmt.Sprintf("LENDER#%d#ID#", lenderID)

// 	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
// 		TableName:              aws.String(r.db.Table),
// 		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":pk":     &types.AttributeValueMemberS{Value: pk},
// 			":prefix": &types.AttributeValueMemberS{Value: prefix},
// 		},
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	var products []models.Product

// 	if err := attributevalue.UnmarshalListOfMaps(out.Items, &products); err != nil {
// 		return nil, err
// 	}

// 	// Convert to product responses
// 	var responses []models.ProductResponse
// 	for _, p := range products {
// 		c, _ := r.categoryRepo.FindByID(p.CategoryID)
// 		u, _ := r.userRepo.FindByID(p.LenderID)

// 		responses = append(responses, models.ProductResponse{
// 			Product:  p,
// 			Category: c,
// 			User:     *u,
// 		})
// 	}

// 	return responses, nil
// }

// func (r *ProductDBRepo) FindByName(name string) ([]models.ProductResponse, error) {
// 	pk := "PRODUCT"
// 	prefix := fmt.Sprintf("NAME#%s", strings.ToLower(name))

// 	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
// 		TableName:              aws.String(r.db.Table),
// 		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":pk":     &types.AttributeValueMemberS{Value: pk},
// 			":prefix": &types.AttributeValueMemberS{Value: prefix},
// 		},
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	var products []models.Product
// 	if err := attributevalue.UnmarshalListOfMaps(out.Items, &products); err != nil {
// 		return nil, err
// 	}

// 	var responses []models.ProductResponse
// 	for _, p := range products {
// 		c, _ := r.categoryRepo.FindByID(p.CategoryID)
// 		u, _ := r.userRepo.FindByID(p.LenderID)

// 		responses = append(responses, models.ProductResponse{
// 			Product:  p,
// 			Category: c,
// 			User:     *u,
// 		})
// 	}

// 	return responses, nil
// }

// func (r *ProductDBRepo) Delete(id int64, lenderID int64, name string) error {
// 	pk := "PRODUCT"

// 	sks := []string{
// 		fmt.Sprintf("PRODUCT#%d", id),
// 		fmt.Sprintf("LENDER#%d#ID#%d", lenderID, id),
// 		fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(name), id),
// 	}

// 	for _, sk := range sks {
// 		_, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
// 			TableName: aws.String(r.db.Table),
// 			Key: map[string]types.AttributeValue{
// 				"pk": &types.AttributeValueMemberS{Value: pk},
// 				"sk": &types.AttributeValueMemberS{Value: sk},
// 			},
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (r *ProductDBRepo) Save() error {
// 	return nil
// }

package product_repo

import (
	"context"
	"errors"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/user_repo"
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

// ✅ Create Product
func (r *ProductDBRepo) Create(product *models.Product) error {
    product.ID = time.Now().UnixNano()
    product.CreatedAt = time.Now()

    // Common attributes
    base := map[string]types.AttributeValue{
        "ID":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.ID)},
        "LenderID":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.LenderID)},
        "CategoryID": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.CategoryID)},
        "Name":       &types.AttributeValueMemberS{Value: product.Name},
        "Description": &types.AttributeValueMemberS{Value: product.Description},
        "Duration":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.Duration)},
        "IsAvailable": &types.AttributeValueMemberBOOL{Value: product.IsAvailable},
        "ImageUrl":   &types.AttributeValueMemberS{Value: product.ImageUrl},
        "CreatedAt":  &types.AttributeValueMemberS{Value: product.CreatedAt.Format(time.RFC3339)},
    }

    // Items for access patterns
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
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s#ID#%d", product.Name, product.ID)},
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

// ✅ FindByID
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

    var product models.Product
    if err := attributevalue.UnmarshalMap(out.Item, &product); err != nil {
        return nil, fmt.Errorf("failed to unmarshal product: %w", err)
    }

    category, _ := r.categoryRepo.FindByID(product.CategoryID)
    user, _ := r.userRepo.FindByID(product.LenderID)

    return &models.ProductResponse{
        Product:  product,
        Category: category,
        User:     *user,
    }, nil
}

// ✅ FindAll (filter by lender or name)
func (r *ProductDBRepo) FindAll(filters models.ProductFilter) ([]*models.ProductResponse, error) {
    var skPrefix string
    if filters.LenderID != "" {
        skPrefix = fmt.Sprintf("LENDER#%s", filters.LenderID)
    } else if filters.Search != "" {
        skPrefix = fmt.Sprintf("NAME#%s", filters.Search)
    } else {
        skPrefix = "PRODUCT#"
    }

    out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
        TableName:              aws.String(r.db.Table),
        KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk":       &types.AttributeValueMemberS{Value: "PRODUCT"},
            ":skPrefix": &types.AttributeValueMemberS{Value: skPrefix},
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to query products: %w", err)
    }

    var products []models.Product
    if err := attributevalue.UnmarshalListOfMaps(out.Items, &products); err != nil {
        return nil, fmt.Errorf("failed to unmarshal products: %w", err)
    }

    var responses []*models.ProductResponse
    for _, p := range products {
        category, _ := r.categoryRepo.FindByID(p.CategoryID)
        user, _ := r.userRepo.FindByID(p.LenderID)
        responses = append(responses, &models.ProductResponse{
            Product:  p,
            Category: category,
            User:     *user,
        })
    }
    return responses, nil
}

// ✅ Update Product
func (r *ProductDBRepo) Update(product *models.Product) error {
    // Delete old copies and recreate
    if err := r.Delete(product.ID); err != nil {
        return err
    }
    return r.Create(product)
}

// ✅ Delete Product
func (r *ProductDBRepo) Delete(id int64) error {
    // Delete all copies by scanning with begins_with(sk, ID#)
    out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
        TableName:              aws.String(r.db.Table),
        KeyConditionExpression: aws.String("pk = :pk AND contains(sk, :id)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
            ":id": &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", id)},
        },
    })
    if err != nil {
        return fmt.Errorf("failed to find product copies: %w", err)
    }

    for _, item := range out.Items {
        var pk, sk string
        attributevalue.Unmarshal(item["pk"], &pk)
        attributevalue.Unmarshal(item["sk"], &sk)
        _, _ = r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
            TableName: aws.String(r.db.Table),
            Key: map[string]types.AttributeValue{
                "pk": &types.AttributeValueMemberS{Value: pk},
                "sk": &types.AttributeValueMemberS{Value: sk},
            },
        })
    }
    return nil
}

// ✅ Save (no-op)
func (r *ProductDBRepo) Save() error {
    return nil
}

// Helper to merge maps
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