package order_repo

import (
	"database/sql"
	"errors"
	"loopit/internal/enums/order_status"
	"loopit/internal/models"
	"loopit/internal/repository/product_repo"
	"time"

	"github.com/lib/pq"
)

type OrderDBRepo struct {
	db          *sql.DB
	productRepo product_repo.ProductRepo
}

func NewOrderDBRepo(db *sql.DB, productRepo product_repo.ProductRepo) *OrderDBRepo {
	return &OrderDBRepo{
		db:          db,
		productRepo: productRepo,
	}
}

// CreateOrder inserts a new order into the database
func (r *OrderDBRepo) CreateOrder(order models.Order) error {
	query := `
	INSERT INTO orders (product_id, user_id, start_date, end_date, total_amount, security_amount, status, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`
	err := r.db.QueryRow(query, order.ProductID, order.UserID, order.StartDate, order.EndDate, order.TotalAmount, order.SecurityAmount, order.Status, time.Now()).Scan(&order.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrderStatus updates the status of an order
func (r *OrderDBRepo) UpdateOrderStatus(orderID int, newStatus string) error {
	statusEnum, err := order_status.ParseStatus(newStatus)
	if err != nil {
		return err
	}

	result, err := r.db.Exec("UPDATE orders SET status=$1 WHERE id=$2", statusEnum, orderID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}

// GetOrderHistory returns orders for a user, optionally filtered by status
func (r *OrderDBRepo) GetOrderHistory(userID int, filterStatuses []string) ([]*models.Order, error) {
	query := "SELECT id, product_id, user_id, start_date, end_date, total_amount, security_amount, status, created_at FROM orders WHERE user_id=$1"
	args := []interface{}{userID}

	if len(filterStatuses) > 0 {
		query += " AND status = ANY($2)"
		args = append(args, pq.Array(filterStatuses)) // requires import "github.com/lib/pq"
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.UserID, &o.StartDate, &o.EndDate, &o.TotalAmount, &o.SecurityAmount, &o.Status, &o.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

// GetLenderOrders returns orders for products owned by a lender
func (r *OrderDBRepo) GetLenderOrders(userID int) ([]*models.Order, error) {
	query := `
	SELECT o.id, o.product_id, o.user_id, o.start_date, o.end_date, o.total_amount, o.security_amount, o.status, o.created_at
	FROM orders o
	JOIN products p ON o.product_id = p.id
	WHERE p.lender_id=$1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.UserID, &o.StartDate, &o.EndDate, &o.TotalAmount, &o.SecurityAmount, &o.Status, &o.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

// GetOrderByID returns a single order by ID
func (r *OrderDBRepo) GetOrderByID(orderID int) (*models.Order, error) {
	row := r.db.QueryRow("SELECT id, product_id, user_id, start_date, end_date, total_amount, security_amount, status, created_at FROM orders WHERE id=$1", orderID)

	var o models.Order
	if err := row.Scan(&o.ID, &o.ProductID, &o.UserID, &o.StartDate, &o.EndDate, &o.TotalAmount, &o.SecurityAmount, &o.Status, &o.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &o, nil
}

// Save is a no-op for Postgres
func (r *OrderDBRepo) Save() error {
	return nil
}
