// package order_repo

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"loopit/internal/db"
// 	"loopit/internal/enums/order_status"
// 	"loopit/internal/models"
// 	"loopit/internal/repository/product_repo"
// 	"time"

// 	"github.com/lib/pq"
// )

// type OrderDBRepo struct {
// 	db          *db.DynamoClient
// 	productRepo product_repo.ProductRepo
// }

// func NewOrderDBRepo(db *db.DynamoClient, productRepo product_repo.ProductRepo) *OrderDBRepo {
// 	return &OrderDBRepo{
// 		db:          db,
// 		productRepo: productRepo,

// 	}
// }

// // CreateOrder inserts a new order into the database
// func (r *OrderDBRepo) CreateOrder(order models.Order) error {
// 	query := `
//     INSERT INTO orders (product_id, user_id, start_date, end_date, total_amount, security_amount, status, created_at)
//     VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
//     RETURNING id
//     `
// 	err := r.db.QueryRow(query, order.ProductID, order.UserID, order.StartDate, order.EndDate, order.TotalAmount, order.SecurityAmount, order.Status.String(), time.Now()).Scan(&order.ID)
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("DB error creating order: %v", err))
// 		return err
// 	}
// 	return nil
// }

// // UpdateOrderStatus updates the status of an order
// func (r *OrderDBRepo) UpdateOrderStatus(orderID int64, newStatus string) error {
// 	var err error
// 	var result sql.Result

// 	if newStatus == "Returned" {
// 		currentTime := time.Now()

// 		result, err = r.db.Exec(
// 			`UPDATE orders SET status=$1, end_date=$2 WHERE id=$3`,
// 			newStatus,
// 			currentTime,
// 			orderID,
// 		)
// 	} else {
// 		result, err = r.db.Exec(
// 			`UPDATE orders SET status=$1 WHERE id=$2`,
// 			newStatus,
// 			orderID,
// 		)
// 	}

// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("DB error updating order %d to status %s: %v", orderID, newStatus, err))
// 		return err
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		r.log.Warning(fmt.Sprintf("DB: No order found to update for id %d", orderID))
// 		return errors.New("order not found")
// 	}

// 	return nil
// }

// // GetOrderHistory returns orders for a user, optionally filtered by status
// func (r *OrderDBRepo) GetOrderHistory(userID int64, filterStatuses []string) ([]*models.Order, error) {
// 	query := "SELECT id, product_id, user_id, start_date, end_date, total_amount, security_amount, status, created_at FROM orders WHERE user_id=$1"
// 	args := []any{userID}

// 	if len(filterStatuses) > 0 {
// 		query += " AND status = ANY($2)"
// 		args = append(args, pq.Array(filterStatuses))
// 	}

// 	rows, err := r.db.Query(query, args...)
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("DB error fetching order history for user %d: %v", userID, err))
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var orders []*models.Order
// 	var statusStr string
// 	for rows.Next() {
// 		var o models.Order
// 		if err := rows.Scan(&o.ID, &o.ProductID, &o.UserID, &o.StartDate, &o.EndDate, &o.TotalAmount, &o.SecurityAmount, &statusStr, &o.CreatedAt); err != nil {
// 			r.log.Warning(fmt.Sprintf("DB warning: could not scan order row: %v", err))
// 			continue
// 		}
// 		o.Status, err = order_status.ParseStatus(statusStr)
// 		if err != nil {
// 			r.log.Warning(fmt.Sprintf("DB warning: could not parse order status: %v", err))
// 			continue
// 		}
// 		orders = append(orders, &o)
// 	}
// 	return orders, nil
// }

// // GetLenderOrders returns orders for products owned by a lender
// func (r *OrderDBRepo) GetLenderOrders(userID int64) ([]*models.Order, error) {
// 	query := `
//     SELECT o.id, o.product_id, o.user_id, o.start_date, o.end_date, o.total_amount, o.security_amount, o.status, o.created_at
//     FROM orders o
//     JOIN products p ON o.product_id = p.id
//     WHERE p.lender_id=$1
//     `
// 	rows, err := r.db.Query(query, userID)
// 	fmt.Println("rows==", *rows)
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("DB error fetching lender orders for user %d: %v", userID, err))
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var orders []*models.Order
// 	var statusStr string
// 	for rows.Next() {
// 		var o models.Order
// 		if err := rows.Scan(&o.ID, &o.ProductID, &o.UserID, &o.StartDate, &o.EndDate, &o.TotalAmount, &o.SecurityAmount, &statusStr, &o.CreatedAt); err != nil {
// 			r.log.Warning(fmt.Sprintf("DB warning: could not scan lender order row: %v", err))
// 			continue
// 		}
// 		o.Status, err = order_status.ParseStatus(statusStr)
// 		if err != nil {
// 			r.log.Warning(fmt.Sprintf("DB warning: could not parse lender order status: %v", err))
// 			continue
// 		}
// 		orders = append(orders, &o)
// 		fmt.Println("orders in repo=", &orders)
// 	}
// 	return orders, nil
// }

// // GetOrderByID returns a single order by ID
// func (r *OrderDBRepo) GetOrderByID(orderID int64) (*models.Order, error) {
// 	row := r.db.QueryRow("SELECT id, product_id, user_id, start_date, end_date, total_amount, security_amount, status, created_at FROM orders WHERE id=$1", orderID)
// 	var o models.Order
// 	var statusStr string
// 	if err := row.Scan(&o.ID, &o.ProductID, &o.UserID, &o.StartDate, &o.EndDate, &o.TotalAmount, &o.SecurityAmount, &statusStr, &o.CreatedAt); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			r.log.Warning(fmt.Sprintf("DB: No order found with id %d", orderID))
// 			return nil, errors.New("order not found")
// 		}
// 		r.log.Error(fmt.Sprintf("DB error fetching order by id %d: %v", orderID, err))
// 		return nil, err
// 	}
// 	status, err := order_status.ParseStatus(statusStr)
// 	if err != nil {
// 		r.log.Warning(fmt.Sprintf("DB warning: could not parse order status for id %d: %v", orderID, err))
// 		return nil, err
// 	}
// 	o.Status = status
// 	return &o, nil
// }

// // Save is a no-op for Postgres
// func (r *OrderDBRepo) Save() error {
// 	return nil
// }

package order_repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"loopit/internal/db"
	"loopit/internal/enums/order_status"
	"loopit/internal/models"
	"loopit/internal/repository/product_repo"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type OrderDBRepo struct {
    db          *db.DynamoClient
    productRepo product_repo.ProductRepo
}

func NewOrderDBRepo(db *db.DynamoClient, productRepo product_repo.ProductRepo) *OrderDBRepo {
    return &OrderDBRepo{
        db:          db,
        productRepo: productRepo,
    }
}

// ✅ CreateOrder
// func (r *OrderDBRepo) CreateOrder(order models.Order) error {
//     order.ID = time.Now().UnixNano()
//     order.CreatedAt = time.Now()

//     // Fetch lender ID from product
//     product, err := r.productRepo.FindByID(order.ProductID)
//     if err != nil {
//         return fmt.Errorf("failed to fetch product for lender info: %w", err)
//     }
//     lenderID := product.Product.LenderID

//     base := map[string]types.AttributeValue{
//         "ID":             &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", order.ID)},
//         "ProductID":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", order.ProductID)},
//         "UserID":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", order.UserID)},
//         "StartDate":      &types.AttributeValueMemberS{Value: order.StartDate.Format(time.RFC3339)},
//         "EndDate":        &types.AttributeValueMemberS{Value: order.EndDate.Format(time.RFC3339)},
//         "TotalAmount":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", order.TotalAmount)},
//         "SecurityAmount": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", order.SecurityAmount)},
//         "Status":         &types.AttributeValueMemberS{Value: order.Status.String()},
//         "CreatedAt":      &types.AttributeValueMemberS{Value: order.CreatedAt.Format(time.RFC3339)},
//     }

//     // Items for access patterns
//     items := []map[string]types.AttributeValue{
//         mergeMap(base, map[string]types.AttributeValue{
//             "pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%d", order.UserID)},
//             "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#ID#%d", order.ID)},
//         }),
//         mergeMap(base, map[string]types.AttributeValue{
//             "pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d", lenderID)},
//             "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#ID#%d", order.ID)},
//         }),
//     }

//     for _, item := range items {
//         _, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
//             TableName: aws.String(r.db.Table),
//             Item:      item,
//         })
//         if err != nil {
//             return fmt.Errorf("failed to create order item: %w", err)
//         }
//     }
//     return nil
// }

func (r *OrderDBRepo) CreateOrder(order models.Order) error {
    order.ID = time.Now().UnixNano()
    order.CreatedAt = time.Now()

    // Fetch lender ID from product
    product, err := r.productRepo.FindByID(order.ProductID)
    if err != nil {
        return fmt.Errorf("failed to fetch product for lender info: %w", err)
    }
    lenderID := product.Product.LenderID

    // Base attributes
    base := map[string]types.AttributeValue{
        "ID":             &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", order.ID)},
        "ProductID":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", order.ProductID)},
        "UserID":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", order.UserID)},
        "StartDate":      &types.AttributeValueMemberS{Value: order.StartDate.Format(time.RFC3339)},
        "EndDate":        &types.AttributeValueMemberS{Value: order.EndDate.Format(time.RFC3339)},
        "TotalAmount":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", order.TotalAmount)},
        "SecurityAmount": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", order.SecurityAmount)},
        "Status":         &types.AttributeValueMemberS{Value: order.Status.String()},
        "CreatedAt":      &types.AttributeValueMemberS{Value: order.CreatedAt.Format(time.RFC3339)},
    }

    // Items for all access patterns
    items := []map[string]types.AttributeValue{
        mergeMap(base, map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%d", order.UserID)},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#ID#%d", order.ID)},
        }),
        mergeMap(base, map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d", lenderID)},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#ID#%d", order.ID)},
        }),
        mergeMap(base, map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "ORDER"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", order.ID)},
        }),
    }

    // Insert all items
    for _, item := range items {
        _, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
            TableName: aws.String(r.db.Table),
            Item:      item,
        })
        if err != nil {
            return fmt.Errorf("failed to create order item: %w", err)
        }
    }

    log.Printf("Order created successfully with ID=%d", order.ID)
    return nil
}


// ✅ UpdateOrderStatus
func (r *OrderDBRepo) UpdateOrderStatus(orderID int64, newStatus string) error {
   
    order, err := r.GetOrderByID(orderID)
    if err != nil {
        return err
    }
    product, _ := r.productRepo.FindByID(order.ProductID)
    lenderID := product.Product.LenderID

    keys := []map[string]types.AttributeValue{
        {"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%d", order.UserID)}, "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#ID#%d", orderID)}},
        {"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d", lenderID)}, "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#ID#%d", orderID)}},
    }

    for _, key := range keys {
        _, err := r.db.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
            TableName: aws.String(r.db.Table),
            Key:       key,
            UpdateExpression: aws.String("SET #s = :status"),
            ExpressionAttributeNames: map[string]string{"#s": "Status"},
            ExpressionAttributeValues: map[string]types.AttributeValue{
                ":status": &types.AttributeValueMemberS{Value: newStatus},
            },
        })
        if err != nil {
            return fmt.Errorf("failed to update order status: %w", err)
        }
    }
    return nil
}

func (r *OrderDBRepo) GetOrderHistory(userID int64, filterStatuses []string) ([]*models.Order, error) {
    out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
        TableName:              aws.String(r.db.Table),
        KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk":       &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%d", userID)},
            ":skPrefix": &types.AttributeValueMemberS{Value: "ORDER#"},
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to query order history: %w", err)
    }

    type orderHelper struct {
        ID             int64     `dynamodbav:"ID"`
        ProductID      int64     `dynamodbav:"ProductID"`
        UserID         int64     `dynamodbav:"UserID"`
        StartDate      time.Time `dynamodbav:"StartDate"`
        EndDate        time.Time `dynamodbav:"EndDate"`
        TotalAmount    float64   `dynamodbav:"TotalAmount"`
        SecurityAmount float64   `dynamodbav:"SecurityAmount"`
        StatusStr      string    `dynamodbav:"Status"`
        CreatedAt      time.Time `dynamodbav:"CreatedAt"`
    }

    var helpers []orderHelper
    if err := attributevalue.UnmarshalListOfMaps(out.Items, &helpers); err != nil {
        return nil, fmt.Errorf("failed to unmarshal orders: %w", err)
    }

    var orders []*models.Order
    for _, h := range helpers {
        status, err := order_status.ParseStatus(h.StatusStr)
        if err != nil {
            continue
        }
        orders = append(orders, &models.Order{
            ID:             h.ID,
            ProductID:      h.ProductID,
            UserID:         h.UserID,
            StartDate:      h.StartDate,
            EndDate:        h.EndDate,
            TotalAmount:    h.TotalAmount,
            SecurityAmount: h.SecurityAmount,
            Status:         status,
            CreatedAt:      h.CreatedAt,
        })
    }

    // Filter by status if provided
    if len(filterStatuses) > 0 {
        var filtered []*models.Order
        statusMap := make(map[string]bool)
        for _, s := range filterStatuses {
            statusMap[s] = true
        }
        for _, o := range orders {
            if statusMap[o.Status.String()] {
                filtered = append(filtered, o)
            }
        }
        return filtered, nil
    }

    return orders, nil
}
func (r *OrderDBRepo) GetLenderOrders(userID int64) ([]*models.Order, error) {
    out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
        TableName:              aws.String(r.db.Table),
        KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk":       &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d", userID)},
            ":skPrefix": &types.AttributeValueMemberS{Value: "ORDER#"},
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to query lender orders: %w", err)
    }

    type orderHelper struct {
        ID             int64     `dynamodbav:"ID"`
        ProductID      int64     `dynamodbav:"ProductID"`
        UserID         int64     `dynamodbav:"UserID"`
        StartDate      time.Time `dynamodbav:"StartDate"`
        EndDate        time.Time `dynamodbav:"EndDate"`
        TotalAmount    float64   `dynamodbav:"TotalAmount"`
        SecurityAmount float64   `dynamodbav:"SecurityAmount"`
        StatusStr      string    `dynamodbav:"Status"`
        CreatedAt      time.Time `dynamodbav:"CreatedAt"`
    }

    var helpers []orderHelper
    if err := attributevalue.UnmarshalListOfMaps(out.Items, &helpers); err != nil {
        return nil, fmt.Errorf("failed to unmarshal orders: %w", err)
    }

    var orders []*models.Order
    for _, h := range helpers {
        status, err := order_status.ParseStatus(h.StatusStr)
        if err != nil {
            continue
        }
        orders = append(orders, &models.Order{
            ID:             h.ID,
            ProductID:      h.ProductID,
            UserID:         h.UserID,
            StartDate:      h.StartDate,
            EndDate:        h.EndDate,
            TotalAmount:    h.TotalAmount,
            SecurityAmount: h.SecurityAmount,
            Status:         status,
            CreatedAt:      h.CreatedAt,
        })
    }

    return orders, nil
}

func (r *OrderDBRepo) GetOrderByID(orderID int64) (*models.Order, error) {
    // Use GetItem with pk=ORDER and sk=ID#<orderId>
    out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
        TableName: aws.String(r.db.Table),
        Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "ORDER"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", orderID)},
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get order: %w", err)
    }
    if out.Item == nil {
        return nil, errors.New("order not found")
    }
    log.Print("item-",out.Item)
    // Helper struct for unmarshalling
    type orderHelper struct {
        ID             int64     `dynamodbav:"ID"`
        ProductID      int64     `dynamodbav:"ProductID"`
        UserID         int64     `dynamodbav:"UserID"`
        StartDate      time.Time `dynamodbav:"StartDate"`
        EndDate        time.Time `dynamodbav:"EndDate"`
        TotalAmount    float64   `dynamodbav:"TotalAmount"`
        SecurityAmount float64   `dynamodbav:"SecurityAmount"`
        StatusStr      string    `dynamodbav:"Status"`
        CreatedAt      time.Time `dynamodbav:"CreatedAt"`
    }

    var helper orderHelper
    if err := attributevalue.UnmarshalMap(out.Item, &helper); err != nil {
        return nil, fmt.Errorf("failed to unmarshal order: %w", err)
    }
    log.Print("helper=",helper)
    // Convert status string to enum
    status, err := order_status.ParseStatus(helper.StatusStr)
    if err != nil {
        return nil, fmt.Errorf("invalid status: %w", err)
    }

    return &models.Order{
        ID:             helper.ID,
        ProductID:      helper.ProductID,
        UserID:         helper.UserID,
        StartDate:      helper.StartDate,
        EndDate:        helper.EndDate,
        TotalAmount:    helper.TotalAmount,
        SecurityAmount: helper.SecurityAmount,
        Status:         status,
        CreatedAt:      helper.CreatedAt,
    }, nil
}


func (r *OrderDBRepo) Save() error {
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