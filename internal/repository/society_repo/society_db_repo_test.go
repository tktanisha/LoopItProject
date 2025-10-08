package society_repo_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/society_repo"
	"loopit/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSocietyDBRepo_Create(t *testing.T) {
	tests := []struct {
		name      string
		society   models.Society
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success - society created",
			society: models.Society{
				Name:     "Green Ville",
				Location: "Downtown",
				Pincode:  "123456",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
	INSERT INTO societies (name, location, pincode)
	VALUES ($1, $2, $3)
	RETURNING id
	`)).
					WithArgs("Green Ville", "Downtown", "123456").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "insert fails",
			society: models.Society{
				Name:     "Blue Ville",
				Location: "Uptown",
				Pincode:  "654321",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
	INSERT INTO societies (name, location, pincode)
	VALUES ($1, $2, $3)
	RETURNING id
	`)).
					WithArgs("Blue Ville", "Uptown", "654321").
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
			repo := society_repo.NewSocietyDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())

			err := repo.Create(tt.society)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocietyDBRepo_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success - societies returned",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, location, pincode FROM societies")).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "location", "pincode",
					}).AddRow(1, "Society A", "City X", "111111"))
			},
			wantErr: false,
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, location, pincode FROM societies")).
					WillReturnError(errors.New("db failure"))
			},
			wantErr: true,
		},
		{
			name: "scan error - skip row",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, location, pincode FROM societies")).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "location", "pincode",
					}).AddRow("invalid", "Society B", "City Y", "222222"))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()

			tt.mockSetup(mock)
			repo := society_repo.NewSocietyDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())

			_, err := repo.FindAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocietyDBRepo_FindByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success - society found",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, location, pincode FROM societies WHERE id=$1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "location", "pincode",
					}).AddRow(1, "Society A", "City X", "111111"))
			},
			wantErr: false,
		},
		{
			name: "no rows",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, location, pincode FROM societies WHERE id=$1")).
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "scan error",
			id:   3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, location, pincode FROM societies WHERE id=$1")).
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "location", "pincode",
					}).AddRow("invalid", "Society B", "City Y", "222222"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()

			tt.mockSetup(mock)
			repo := society_repo.NewSocietyDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())

			_, err := repo.FindByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocietyDBRepo_Save(t *testing.T) {
	repo := society_repo.NewSocietyDBRepo(nil, nil)
	if err := repo.Save(); err != nil {
		t.Errorf("Save() expected no error, got %v", err)
	}
}
