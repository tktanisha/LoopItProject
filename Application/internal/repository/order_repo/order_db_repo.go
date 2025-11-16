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
	"fmt"
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

func (r *OrderDBRepo) CreateOrder(o models.Order) error {
	o.ID = time.Now().UnixNano()
	o.CreatedAt = time.Now()

	table := r.db.Table

	// lender fetch
	p, err := r.productRepo.FindByID(o.ProductID)
	if err != nil {
		return fmt.Errorf("could not find lender for product")
	}
	lenderID := p.Product.LenderID

	skID := fmt.Sprintf("ID#%d", o.ID)
	skUser := fmt.Sprintf("ORDER#ID#%d", o.ID)
	skLender := fmt.Sprintf("ORDER#ID#%d", o.ID)

	item := map[string]interface{}{
		"ID":             o.ID,
		"ProductID":      o.ProductID,
		"UserID":         o.UserID,
		"StartDate":      o.StartDate.Format(time.RFC3339),
		"EndDate":        o.EndDate.Format(time.RFC3339),
		"TotalAmount":    o.TotalAmount,
		"SecurityAmount": o.SecurityAmount,
		"Status":         o.Status.String(),
		"CreatedAt":      o.CreatedAt.Format(time.RFC3339),
	}

	entries := []struct{ PK, SK string }{
		{"ORDER", skID},
		{fmt.Sprintf("USER#%d", o.UserID), skUser},
		{fmt.Sprintf("LENDER#%d", lenderID), skLender},
	}

	for _, e := range entries {

		keys, err := attributevalue.MarshalMap(map[string]string{
			"pk": e.PK,
			"sk": e.SK,
		})
		if err != nil {
			return err
		}

		body, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}

		for k, v := range body {
			keys[k] = v
		}

		_, err = r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item:      keys,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderDBRepo) GetOrderByID(id int64) (*models.Order, error) {
	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "ORDER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
		},
	})

	if err != nil || out.Item == nil {
		return nil, fmt.Errorf("order not found")
	}

	var o models.Order
	if err := attributevalue.UnmarshalMap(out.Item, &o); err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *OrderDBRepo) GetOrderHistory(userID int64, statuses []string) ([]models.Order, error) {
	pk := fmt.Sprintf("USER#%d", userID)

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :p)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
			":p":  &types.AttributeValueMemberS{Value: "ORDER#ID#"},
		},
	})

	if err != nil {
		return nil, err
	}

	var orders []models.Order
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &orders); err != nil {
		return nil, err
	}

	// status filter
	if len(statuses) > 0 {
		set := map[string]bool{}
		for _, s := range statuses {
			set[s] = true
		}

		tmp := []models.Order{}
		for _, o := range orders {
			if set[o.Status.String()] {
				tmp = append(tmp, o)
			}
		}
		return tmp, nil
	}

	return orders, nil
}

func (r *OrderDBRepo) GetLenderOrders(lenderID int64) ([]models.Order, error) {
	pk := fmt.Sprintf("LENDER#%d", lenderID)

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :p)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
			":p":  &types.AttributeValueMemberS{Value: "ORDER#ID#"},
		},
	})

	if err != nil {
		return nil, err
	}

	var orders []models.Order
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderDBRepo) UpdateOrderStatus(id int64, newStatus string) error {
	o, err := r.GetOrderByID(id)
	if err != nil {
		return err
	}

	p, err := r.productRepo.FindByID(o.ProductID)
	if err != nil {
		return err
	}

	// remove old items
	r.DeleteOrder(id, o.UserID, p.Product.LenderID)

	// modify
	o.Status, _ = order_status.ParseStatus(newStatus)

	// recreate all rows
	return r.CreateOrder(o)
}

func (r *OrderDBRepo) DeleteOrder(id, userID, lenderID int64) error {
	sks := []string{
		fmt.Sprintf("ID#%d", id),
		fmt.Sprintf("ORDER#ID#%d", id),
		fmt.Sprintf("ORDER#ID#%d", id),
	}
	pks := []string{
		"ORDER",
		fmt.Sprintf("USER#%d", userID),
		fmt.Sprintf("LENDER#%d", lenderID),
	}

	for i := range sks {
		_, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String(r.db.Table),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: pks[i]},
				"sk": &types.AttributeValueMemberS{Value: sks[i]},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderDBRepo) Save() error { return nil }
