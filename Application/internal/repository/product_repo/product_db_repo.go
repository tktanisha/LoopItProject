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

    // Common attributes
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

    // Items for different access patterns
    items := []map[string]types.AttributeValue{
        // Primary item by product ID
        mergeMap(base, map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)},
        }),
        // Item by lender
        mergeMap(base, map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d#ID#%d", product.LenderID, product.ID)},
        }),
        // Item by name
        mergeMap(base, map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(product.Name), product.ID)},
            
        }),
        // ✅ Item by category for filtering
        mergeMap(base, map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("CATEGORY#%d", product.CategoryID)},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)},
        }),
    }

    // Write all items to DynamoDB
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

    // Helper struct for unmarshalling
    type productHelper struct {
        ID          int64   `dynamodbav:"ID"`
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

    // Convert string IDs to int64
    lenderID, _ := strconv.ParseInt(helper.LenderID, 10, 64)
    categoryID, _ := strconv.ParseInt(helper.CategoryID, 10, 64)

    // Fetch related category
    var category models.Category
    if r.categoryRepo != nil {
        category, _ = r.categoryRepo.FindByID(categoryID)
    }

    // Fetch related user
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

    // Helper struct for unmarshalling
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

    // ✅ Apply in-memory search filter for name or description
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

    // ✅ Apply IsAvailable filter if provided
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

// func (r *ProductDBRepo) Update(product *models.Product) error {
//     ctx := context.TODO()

//     // Common update expression
//     updateExpr := "SET #n = :name, Description = :description, #dur = :duration, IsAvailable = :isAvailable, CategoryID = :categoryId, LenderID = :lenderId"

//     exprAttrNames := map[string]string{
//         "#n": "Name", 
//         "#dur": "Duration",
//     }

//     exprAttrValues := map[string]types.AttributeValue{
//         ":name":        &types.AttributeValueMemberS{Value: product.Name},
//         ":description": &types.AttributeValueMemberS{Value: product.Description},
//         ":duration":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", product.Duration)},
//         ":isAvailable": &types.AttributeValueMemberBOOL{Value: product.IsAvailable},
//         ":categoryId":  &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", product.CategoryID)},
//         ":lenderId":    &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", product.LenderID)},
//     }

//     _, err := r.db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
//         TableName: aws.String(r.db.Table),
//         Key: map[string]types.AttributeValue{
//             "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
//             "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)},
//         },
//         UpdateExpression:          aws.String(updateExpr),
//         ExpressionAttributeNames:  exprAttrNames,
//         ExpressionAttributeValues: exprAttrValues,
//     })
//     if err != nil {
//         return fmt.Errorf("failed to update main product record: %w", err)
//     }

//     _, err = r.db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
//         TableName: aws.String(r.db.Table),
//         Key: map[string]types.AttributeValue{
//             "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
//             "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d#ID#%d", product.LenderID, product.ID)},
//         },
//         UpdateExpression:          aws.String(updateExpr),
//         ExpressionAttributeNames:  exprAttrNames,
//         ExpressionAttributeValues: exprAttrValues,
//     })
//     if err != nil {
//         return fmt.Errorf("failed to update lender index record: %w", err)
//     }

//     _, err = r.db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
//         TableName: aws.String(r.db.Table),
//         Key: map[string]types.AttributeValue{
//             "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
//             "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s#ID#%d", product.Name, product.ID)},
//         },
//         UpdateExpression:          aws.String(updateExpr),
//         ExpressionAttributeNames:  exprAttrNames,
//         ExpressionAttributeValues: exprAttrValues,
//     })
//     if err != nil {
//         return fmt.Errorf("failed to update name index record: %w", err)
//     }

//     return nil
// }

func (r *ProductDBRepo) Update(product *models.Product) error {
    ctx := context.TODO()

    // 1. Fetch existing product to get old CategoryID and Name
    existing, err := r.FindByID(product.ID)
    if err != nil {
        return fmt.Errorf("failed to fetch existing product: %w", err)
    }

    // 2. Delete old category index if category changed
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

    // 3. Delete old name index if name changed
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

    // 4. Update all current keys (main, lender, new name, new category)
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
            TableName: aws.String(r.db.Table),
            Key:       key,
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

    // 1. Fetch product details
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

    // Unmarshal helper
    type productHelper struct {
        Name       string `dynamodbav:"Name"`
        LenderID   int64  `dynamodbav:"LenderID"`
        CategoryID int64  `dynamodbav:"CategoryID"`
    }
    var p productHelper
    if err := attributevalue.UnmarshalMap(out.Item, &p); err != nil {
        return fmt.Errorf("failed to unmarshal product: %w", err)
    }

    // 2. Prepare all keys to delete
    keys := []map[string]types.AttributeValue{
        { // Main product record
            "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", productID)},
        },
        { // Lender index
            "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d#ID#%d", p.LenderID, productID)},
        },
        { // Name index
            "pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s#ID#%d", strings.ToLower(p.Name), productID)},
        },
        { // Category index
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