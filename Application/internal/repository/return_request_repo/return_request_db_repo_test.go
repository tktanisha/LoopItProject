package return_request_repo_test

// import (
// 	"database/sql"
// 	"errors"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"loopit/internal/db"
// 	"loopit/internal/enums/return_request_status"
// 	"loopit/internal/models"
// 	"loopit/internal/repository/return_request_repo"
// 	"loopit/pkg/logger"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/lib/pq"
// )

// func TestReturnRequestDBRepo_CreateReturnRequest(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		req       models.ReturnRequest
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success - return request created",
// 			req:  models.ReturnRequest{OrderID: 10, Status: return_request_status.Pending},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(`
//     INSERT INTO return_requests (order_id, status, created_at)
//     VALUES ($1, $2, $3)
//     RETURNING id
//     `)).
// 					WithArgs(10, "Pending", sqlmock.AnyArg()).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "insert fails",
// 			req:  models.ReturnRequest{OrderID: 20, Status: return_request_status.Approved},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(`
//     INSERT INTO return_requests (order_id, status, created_at)
//     VALUES ($1, $2, $3)
//     RETURNING id
//     `)).
// 					WithArgs(20, "Approved", sqlmock.AnyArg()).
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
// 			repo := return_request_repo.NewReturnRequestDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())

// 			err := repo.CreateReturnRequest(tt.req)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CreateReturnRequest() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestReturnRequestDBRepo_UpdateReturnRequestStatus(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		id        int
// 		status    string
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name:   "success - status updated",
// 			id:     1,
// 			status: "Approved",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta("UPDATE return_requests SET status=$1 WHERE id=$2")).
// 					WithArgs("Approved", 1).
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:   "update fails - db error",
// 			id:     2,
// 			status: "Rejected",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta("UPDATE return_requests SET status=$1 WHERE id=$2")).
// 					WithArgs("Rejected", 2).
// 					WillReturnError(errors.New("db failure"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:   "no rows affected",
// 			id:     3,
// 			status: "Pending",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta("UPDATE return_requests SET status=$1 WHERE id=$2")).
// 					WithArgs("Pending", 3).
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
// 			repo := return_request_repo.NewReturnRequestDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())

// 			err := repo.UpdateReturnRequestStatus(tt.id, tt.status)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UpdateReturnRequestStatus() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestReturnRequestDBRepo_GetAllReturnRequests(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		filters   []string
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name:    "success - with filter",
// 			filters: []string{"Pending"},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests WHERE status = ANY($1)")).
// 					WithArgs(pq.Array([]string{"Pending"})).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "status", "created_at"}).
// 						AddRow(1, 101, "Pending", time.Now()))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "success - no filter",
// 			filters: nil,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests")).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "status", "created_at"}).
// 						AddRow(2, 102, "Approved", time.Now()))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "query fails",
// 			filters: nil,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests")).
// 					WillReturnError(errors.New("db query failed"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:    "scan error - skipped row",
// 			filters: nil,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests")).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "status", "created_at"}).
// 						AddRow("invalid", 103, "Pending", time.Now()))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "invalid status - skipped row",
// 			filters: nil,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests")).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "status", "created_at"}).
// 						AddRow(3, 104, "unknown_status", time.Now()))
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()

// 			tt.mockSetup(mock)
// 			repo := return_request_repo.NewReturnRequestDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())

// 			_, err := repo.GetAllReturnRequests(tt.filters)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetAllReturnRequests() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestReturnRequestDBRepo_GetReturnRequestByID(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		id        int
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success - found",
// 			id:   1,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests WHERE id=$1")).
// 					WithArgs(1).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "status", "created_at"}).
// 						AddRow(1, 201, "Pending", time.Now()))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "not found",
// 			id:   2,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests WHERE id=$1")).
// 					WithArgs(2).
// 					WillReturnError(sql.ErrNoRows)
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "scan error",
// 			id:   3,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests WHERE id=$1")).
// 					WithArgs(3).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "status", "created_at"}).
// 						AddRow("invalid", 202, "Pending", time.Now()))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid status parsing",
// 			id:   4,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, status, created_at FROM return_requests WHERE id=$1")).
// 					WithArgs(4).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "status", "created_at"}).
// 						AddRow(4, 203, "unknown_status", time.Now()))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()

// 			tt.mockSetup(mock)
// 			repo := return_request_repo.NewReturnRequestDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())

// 			_, err := repo.GetReturnRequestByID(tt.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetReturnRequestByID() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestReturnRequestDBRepo_Save(t *testing.T) {
// 	repo := return_request_repo.NewReturnRequestDBRepo(nil, nil)
// 	if err := repo.Save(); err != nil {
// 		t.Errorf("Save() expected no error, got %v", err)
// 	}
// }
