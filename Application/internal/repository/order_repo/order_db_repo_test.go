package order_repo_test

// import (
// 	"database/sql"
// 	"errors"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"loopit/internal/db"
// 	"loopit/internal/enums/order_status"
// 	custom_mock "loopit/internal/mock"
// 	"loopit/internal/models"
// 	"loopit/internal/repository/order_repo"
// 	"loopit/pkg/logger"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/golang/mock/gomock"
// )

// func TestOrderDBRepo_CreateOrder(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		order     models.Order
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success - order created",
// 			order: models.Order{
// 				ProductID:      1,
// 				UserID:         10,
// 				StartDate:      time.Now(),
// 				EndDate:        time.Now().Add(24 * time.Hour),
// 				TotalAmount:    150.0,
// 				SecurityAmount: 50.0,
// 				Status:         order_status.InUse,
// 			},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(`
//     INSERT INTO orders (product_id, user_id, start_date, end_date, total_amount, security_amount, status, created_at)
//     VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
//     RETURNING id
//     `)).
// 					WithArgs(1, 10, sqlmock.AnyArg(), sqlmock.AnyArg(), 150.0, 50.0, order_status.InUse.String(), sqlmock.AnyArg()).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "failure - db insert error",
// 			order: models.Order{
// 				ProductID:      2,
// 				UserID:         20,
// 				StartDate:      time.Now(),
// 				EndDate:        time.Now().Add(24 * time.Hour),
// 				TotalAmount:    200.0,
// 				SecurityAmount: 60.0,
// 				Status:         order_status.InUse,
// 			},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(`
//     INSERT INTO orders (product_id, user_id, start_date, end_date, total_amount, security_amount, status, created_at)
//     VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
//     RETURNING id
//     `)).
// 					WithArgs(2, 20, sqlmock.AnyArg(), sqlmock.AnyArg(), 200.0, 60.0, order_status.InUse.String(), sqlmock.AnyArg()).
// 					WillReturnError(errors.New("insert failed"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			repo := order_repo.NewOrderDBRepo(db.NewRealDatabase(sqlDB), custom_mock.NewMockProductRepo(ctrl), logger.NewFakeLogger())
// 			err := repo.CreateOrder(tt.order)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CreateOrder() expected error=%v, got %v", tt.wantErr, err)
// 			}
// 		})
// 	}
// }

// func TestOrderDBRepo_UpdateOrderStatus(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		orderID   int
// 		status    string
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name:    "success - status updated",
// 			orderID: 1,
// 			status:  "completed",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status=$1 WHERE id=$2")).
// 					WithArgs("completed", 1).
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "failure - db exec error",
// 			orderID: 2,
// 			status:  "cancelled",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status=$1 WHERE id=$2")).
// 					WithArgs("cancelled", 2).
// 					WillReturnError(errors.New("update failed"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:    "failure - no rows affected",
// 			orderID: 3,
// 			status:  "pending",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status=$1 WHERE id=$2")).
// 					WithArgs("pending", 3).
// 					WillReturnResult(sqlmock.NewResult(0, 0))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			repo := order_repo.NewOrderDBRepo(db.NewRealDatabase(sqlDB), custom_mock.NewMockProductRepo(ctrl), logger.NewFakeLogger())
// 			err := repo.UpdateOrderStatus(tt.orderID, tt.status)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UpdateOrderStatus() expected error=%v, got %v", tt.wantErr, err)
// 			}
// 		})
// 	}
// }

// func TestOrderDBRepo_GetOrderHistory(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		userID         int
// 		filterStatuses []string
// 		mockSetup      func(sqlmock.Sqlmock)
// 		wantErr        bool
// 		expectEmpty    bool
// 	}{
// 		{
// 			name:   "success - orders returned",
// 			userID: 1,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "start_date", "end_date", "total_amount", "security_amount", "status", "created_at"}).
// 					AddRow(1, 101, 1, time.Now(), time.Now().Add(24*time.Hour), 200.0, 50.0, order_status.Returned.String(), time.Now())
// 				mock.ExpectQuery("SELECT id, product_id, user_id.*FROM orders WHERE user_id=.*").
// 					WithArgs(1).
// 					WillReturnRows(rows)
// 			},
// 		},
// 		{
// 			name:   "failure - db query error",
// 			userID: 2,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery("SELECT id, product_id, user_id.*FROM orders WHERE user_id=.*").
// 					WithArgs(2).
// 					WillReturnError(errors.New("db error"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:   "warning - scan failure skips row",
// 			userID: 3,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "start_date", "end_date", "total_amount", "security_amount", "status", "created_at"}).
// 					AddRow("invalid", 101, 3, time.Now(), time.Now().Add(24*time.Hour), 200.0, 50.0, order_status.Returned.String(), time.Now())
// 				mock.ExpectQuery("SELECT id, product_id, user_id.*FROM orders WHERE user_id=.*").
// 					WithArgs(3).
// 					WillReturnRows(rows)
// 			},
// 			expectEmpty: true,
// 		},
// 		{
// 			name:   "warning - status parse failure skips row",
// 			userID: 4,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "start_date", "end_date", "total_amount", "security_amount", "status", "created_at"}).
// 					AddRow(4, 101, 4, time.Now(), time.Now().Add(24*time.Hour), 200.0, 50.0, "invalid_status", time.Now())
// 				mock.ExpectQuery("SELECT id, product_id, user_id.*FROM orders WHERE user_id=.*").
// 					WithArgs(4).
// 					WillReturnRows(rows)
// 			},
// 			expectEmpty: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			repo := order_repo.NewOrderDBRepo(db.NewRealDatabase(sqlDB), custom_mock.NewMockProductRepo(ctrl), logger.NewFakeLogger())
// 			orders, err := repo.GetOrderHistory(tt.userID, tt.filterStatuses)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetOrderHistory() expected error=%v, got %v", tt.wantErr, err)
// 			}
// 			if tt.expectEmpty && len(orders) != 0 {
// 				t.Errorf("expected empty orders but got %d", len(orders))
// 			}
// 		})
// 	}
// }

// func TestOrderDBRepo_GetLenderOrders(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		userID      int
// 		mockSetup   func(sqlmock.Sqlmock)
// 		wantErr     bool
// 		expectEmpty bool
// 	}{
// 		{
// 			name:   "success - lender orders returned",
// 			userID: 1,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "start_date", "end_date", "total_amount", "security_amount", "status", "created_at"}).
// 					AddRow(1, 101, 10, time.Now(), time.Now().Add(24*time.Hour), 300.0, 100.0, order_status.Returned.String(), time.Now())
// 				mock.ExpectQuery("SELECT o.id, o.product_id.*FROM orders o.*").
// 					WithArgs(1).
// 					WillReturnRows(rows)
// 			},
// 		},
// 		{
// 			name:   "failure - db error",
// 			userID: 2,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery("SELECT o.id, o.product_id.*FROM orders o.*").
// 					WithArgs(2).
// 					WillReturnError(errors.New("db error"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:   "warning - scan failure skips row",
// 			userID: 3,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "start_date", "end_date", "total_amount", "security_amount", "status", "created_at"}).
// 					AddRow("bad", 101, 10, time.Now(), time.Now().Add(24*time.Hour), 300.0, 100.0, order_status.Returned.String(), time.Now())
// 				mock.ExpectQuery("SELECT o.id, o.product_id.*FROM orders o.*").
// 					WithArgs(3).
// 					WillReturnRows(rows)
// 			},
// 			expectEmpty: true,
// 		},
// 		{
// 			name:   "warning - status parse failure skips row",
// 			userID: 4,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "start_date", "end_date", "total_amount", "security_amount", "status", "created_at"}).
// 					AddRow(4, 101, 10, time.Now(), time.Now().Add(24*time.Hour), 300.0, 100.0, "bad_status", time.Now())
// 				mock.ExpectQuery("SELECT o.id, o.product_id.*FROM orders o.*").
// 					WithArgs(4).
// 					WillReturnRows(rows)
// 			},
// 			expectEmpty: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			repo := order_repo.NewOrderDBRepo(db.NewRealDatabase(sqlDB), custom_mock.NewMockProductRepo(ctrl), logger.NewFakeLogger())
// 			orders, err := repo.GetLenderOrders(tt.userID)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetLenderOrders() expected error=%v, got %v", tt.wantErr, err)
// 			}
// 			if tt.expectEmpty && len(orders) != 0 {
// 				t.Errorf("expected empty orders but got %d", len(orders))
// 			}
// 		})
// 	}
// }

// func TestOrderDBRepo_GetOrderByID(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		orderID   int
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name:    "success - order returned",
// 			orderID: 1,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "start_date", "end_date", "total_amount", "security_amount", "status", "created_at"}).
// 					AddRow(1, 101, 10, time.Now(), time.Now().Add(24*time.Hour), 250.0, 75.0, order_status.InUse.String(), time.Now())
// 				mock.ExpectQuery("SELECT id, product_id, user_id.*FROM orders WHERE id=.*").
// 					WithArgs(1).
// 					WillReturnRows(rows)
// 			},
// 		},
// 		{
// 			name:    "failure - no rows",
// 			orderID: 2,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery("SELECT id, product_id, user_id.*FROM orders WHERE id=.*").
// 					WithArgs(2).
// 					WillReturnError(sql.ErrNoRows)
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:    "failure - db error",
// 			orderID: 3,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery("SELECT id, product_id, user_id.*FROM orders WHERE id=.*").
// 					WithArgs(3).
// 					WillReturnError(errors.New("db failure"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:    "failure - status parse error",
// 			orderID: 4,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "start_date", "end_date", "total_amount", "security_amount", "status", "created_at"}).
// 					AddRow(4, 101, 10, time.Now(), time.Now().Add(24*time.Hour), 250.0, 75.0, "bad_status", time.Now())
// 				mock.ExpectQuery("SELECT id, product_id, user_id.*FROM orders WHERE id=.*").
// 					WithArgs(4).
// 					WillReturnRows(rows)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			repo := order_repo.NewOrderDBRepo(db.NewRealDatabase(sqlDB), custom_mock.NewMockProductRepo(ctrl), logger.NewFakeLogger())
// 			_, err := repo.GetOrderByID(tt.orderID)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetOrderByID() expected error=%v, got %v", tt.wantErr, err)
// 			}
// 		})
// 	}
// }

// func TestOrderDBRepo_Save(t *testing.T) {
// 	sqlDB, _, _ := sqlmock.New()
// 	defer sqlDB.Close()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	repo := order_repo.NewOrderDBRepo(db.NewRealDatabase(sqlDB), custom_mock.NewMockProductRepo(ctrl), logger.NewFakeLogger())
// 	if err := repo.Save(); err != nil {
// 		t.Errorf("Save() expected nil, got %v", err)
// 	}
// }
