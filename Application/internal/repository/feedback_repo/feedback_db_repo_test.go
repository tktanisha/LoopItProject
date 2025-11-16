package feedback_repo_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/feedback_repo"
	"loopit/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFeedBackDBRepo_CreateFeedback(t *testing.T) {
	tests := []struct {
		name      string
		feedback  models.Feedback
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success - feedback created",
			feedback: models.Feedback{
				GivenBy: 1, GivenTo: 2, Text: "Great service", Rating: 5,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
    INSERT INTO feedbacks (given_by, given_to, text, rating, created_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id
    `)).
					WithArgs(1, 2, "Great service", 5, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "db error - insert fails",
			feedback: models.Feedback{
				GivenBy: 10, GivenTo: 20, Text: "Failed feedback", Rating: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
    INSERT INTO feedbacks (given_by, given_to, text, rating, created_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id
    `)).
					WithArgs(10, 20, "Failed feedback", 1, sqlmock.AnyArg()).
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()
			tt.mockSetup(mock)

			repo := feedback_repo.NewFeedBackDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
			err := repo.CreateFeedback(tt.feedback)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestFeedBackDBRepo_GetAllFeedbacks(t *testing.T) {
	currentTime := time.Now()
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		wantNil   bool
		wantCount int
	}{
		{
			name: "success - multiple feedbacks",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "given_by", "given_to", "text", "rating", "created_at"}).
					AddRow(1, 10, 20, "Good", 5, currentTime).
					AddRow(2, 11, 21, "Average", 3, currentTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, given_by, given_to, text, rating, created_at FROM feedbacks")).
					WillReturnRows(rows)
			},
			wantErr:   false,
			wantNil:   false,
			wantCount: 2,
		},
		{
			name: "db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, given_by, given_to, text, rating, created_at FROM feedbacks")).
					WillReturnError(errors.New("db failure"))
			},
			wantErr:   true,
			wantNil:   true,
			wantCount: 0,
		},
		{
			name: "scan error on one row (skipped)",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "given_by", "given_to", "text", "rating", "created_at"}).
					AddRow(nil, 10, 20, "Broken", 1, currentTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, given_by, given_to, text, rating, created_at FROM feedbacks")).
					WillReturnRows(rows)
			},
			wantErr:   false,
			wantNil:   true, // because no valid rows
			wantCount: 0,
		},
		{
			name: "no rows returned",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "given_by", "given_to", "text", "rating", "created_at"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, given_by, given_to, text, rating, created_at FROM feedbacks")).
					WillReturnRows(rows)
			},
			wantErr:   false,
			wantNil:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()
			tt.mockSetup(mock)

			repo := feedback_repo.NewFeedBackDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
			feedbacks, err := repo.GetAllFeedbacks()

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}

			if tt.wantNil && feedbacks != nil {
				t.Errorf("expected nil feedbacks, got %v", feedbacks)
			}

			if feedbacks != nil && len(feedbacks) != tt.wantCount {
				t.Errorf("expected %d feedbacks, got %d", tt.wantCount, len(feedbacks))
			}
		})
	}
}

func TestFeedBackDBRepo_Save(t *testing.T) {
	sqlDB, _, _ := sqlmock.New()
	defer sqlDB.Close()

	repo := feedback_repo.NewFeedBackDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
	if err := repo.Save(); err != nil {
		t.Errorf("expected no error from Save(), got %v", err)
	}
}
