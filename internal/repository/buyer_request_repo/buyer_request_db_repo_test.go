package buyer_request_repo_test

// import (
// 	"database/sql"
// 	"errors"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"loopit/internal/db"
// 	"loopit/internal/enums/buyer_request_status"
// 	"loopit/internal/models"
// 	"loopit/internal/repository/buyer_request_repo"
// 	"loopit/pkg/logger"

// 	"github.com/DATA-DOG/go-sqlmock"
// )

// func TestBuyerRequestDBRepo_GetAllBuyerRequests(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		filter        []string
// 		mockSetup     func(mock sqlmock.Sqlmock)
// 		wantErr       bool
// 		expectedCount int
// 	}{
// 		{
// 			name:   "success - no filter",
// 			filter: nil,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests")).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "requested_by", "status", "created_at"}).
// 						AddRow(1, 100, 200, buyer_request_status.Pending.String(), time.Now()))
// 			},
// 			wantErr:       false,
// 			expectedCount: 1,
// 		},
// 		{
// 			name:   "success - with filter",
// 			filter: []string{buyer_request_status.Approved.String()},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests WHERE status = ANY($1)")).
// 					WithArgs(sqlmock.AnyArg()).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "requested_by", "status", "created_at"}).
// 						AddRow(2, 101, 201, buyer_request_status.Approved.String(), time.Now()))
// 			},
// 			wantErr:       false,
// 			expectedCount: 1,
// 		},
// 		{
// 			name:   "query error",
// 			filter: nil,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests")).
// 					WillReturnError(errors.New("db error"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:   "scan error - continue",
// 			filter: nil,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests")).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "requested_by", "status", "created_at"}).
// 						AddRow("bad_id", 100, 200, buyer_request_status.Pending.String(), time.Now()))
// 			},
// 			wantErr:       false,
// 			expectedCount: 0,
// 		},
// 		{
// 			name:   "invalid status - skip row",
// 			filter: nil,
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests")).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "requested_by", "status", "created_at"}).
// 						AddRow(1, 100, 200, "INVALID_STATUS", time.Now()))
// 			},
// 			wantErr:       false,
// 			expectedCount: 0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			repo := buyer_request_repo.NewBuyerRequestDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
// 			got, err := repo.GetAllBuyerRequests(tt.filter)

// 			if tt.wantErr && err == nil {
// 				t.Errorf("expected error, got none")
// 			}
// 			if !tt.wantErr && len(got) != tt.expectedCount {
// 				t.Errorf("expected %d results, got %d", tt.expectedCount, len(got))
// 			}
// 		})
// 	}
// }

// func TestBuyerRequestDBRepo_UpdateStatusBuyerRequest(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		mockSetup func(mock sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta(
// 					"UPDATE buying_requests SET status=$1 WHERE id=$2")).
// 					WithArgs(buyer_request_status.Approved.String(), 1).
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "no rows affected",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta(
// 					"UPDATE buying_requests SET status=$1 WHERE id=$2")).
// 					WithArgs(buyer_request_status.Approved.String(), 1).
// 					WillReturnResult(sqlmock.NewResult(1, 0))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "query error",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec(regexp.QuoteMeta(
// 					"UPDATE buying_requests SET status=$1 WHERE id=$2")).
// 					WithArgs(buyer_request_status.Approved.String(), 1).
// 					WillReturnError(errors.New("db error"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			repo := buyer_request_repo.NewBuyerRequestDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
// 			err := repo.UpdateStatusBuyerRequest(1, buyer_request_status.Approved.String())

// 			if tt.wantErr && err == nil {
// 				t.Errorf("expected error, got none")
// 			}
// 		})
// 	}
// }

// func TestBuyerRequestDBRepo_CreateBuyerRequest(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		mockSetup func(mock sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(`
// 	INSERT INTO buying_requests (product_id, requested_by, status, created_at)
// 	VALUES ($1, $2, $3, $4)
// 	RETURNING id
// 	`)).
// 					WithArgs(100, 200, buyer_request_status.Pending.String(), sqlmock.AnyArg()).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "query error",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(`
// 	INSERT INTO buying_requests (product_id, requested_by, status, created_at)
// 	VALUES ($1, $2, $3, $4)
// 	RETURNING id
// 	`)).
// 					WithArgs(100, 200, buyer_request_status.Pending.String(), sqlmock.AnyArg()).
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

// 			repo := buyer_request_repo.NewBuyerRequestDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
// 			err := repo.CreateBuyerRequest(models.BuyingRequest{
// 				ProductID:   100,
// 				RequestedBy: 200,
// 				Status:      buyer_request_status.Pending,
// 			})

// 			if tt.wantErr && err == nil {
// 				t.Errorf("expected error, got none")
// 			}
// 		})
// 	}
// }

// func TestBuyerRequestDBRepo_GetBuyerRequestByID(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		mockSetup func(mock sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests WHERE id=$1")).
// 					WithArgs(1).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "requested_by", "status", "created_at"}).
// 						AddRow(1, 100, 200, buyer_request_status.Pending.String(), time.Now()))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "not found",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests WHERE id=$1")).
// 					WithArgs(1).
// 					WillReturnError(sql.ErrNoRows)
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "query error",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests WHERE id=$1")).
// 					WithArgs(1).
// 					WillReturnError(errors.New("db error"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid status",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(
// 					"SELECT id, product_id, requested_by, status, created_at FROM buying_requests WHERE id=$1")).
// 					WithArgs(1).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "requested_by", "status", "created_at"}).
// 						AddRow(1, 100, 200, "INVALID", time.Now()))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			repo := buyer_request_repo.NewBuyerRequestDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
// 			_, err := repo.GetBuyerRequestByID(1)

// 			if tt.wantErr && err == nil {
// 				t.Errorf("expected error, got none")
// 			}
// 		})
// 	}
// }

// func TestBuyerRequestDBRepo_Save(t *testing.T) {
// 	repo := buyer_request_repo.NewBuyerRequestDBRepo(nil, logger.NewFakeLogger())
// 	if err := repo.Save(); err != nil {
// 		t.Errorf("Save() expected nil, got %v", err)
// 	}
// }
