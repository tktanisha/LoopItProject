package initializer_test

import (
	"testing"

	"loopit/internal/config"
	"loopit/internal/db"
	"loopit/internal/initializer"
	"loopit/internal/mock"
	"loopit/pkg/logger"

	"github.com/golang/mock/gomock"
)

func mockInitDBRepos(t *testing.T, fn func(logger.LoggerInterface, *db.DynamoClient) error) func() {
	original := initializer.InitDBRepos
	initializer.InitDBRepos = func(l logger.LoggerInterface, d *db.DynamoClient) error {
		return fn(l, d)
	}
	return func() { initializer.InitDBRepos = original }
}

func TestInitServices_DBStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	// Force config to "db"
	originalStorageType := config.AppConfig.StorageType
	config.AppConfig.StorageType = "db"
	defer func() { config.AppConfig.StorageType = originalStorageType }()

	// Mock InitDBRepos to avoid side effects
	restore := mockInitDBRepos(t, func(l logger.LoggerInterface, db *db.DynamoClient) error {
		if db != mockDB {
			t.Errorf("InitDBRepos() received unexpected db instance")
		}
		if l == nil {
			t.Errorf("InitDBRepos() received nil logger")
		}
		return nil
	})
	defer restore()

	// Call function under test
	log := logger.NewFakeLogger()
	err := initializer.InitServices(mockDB, log)
	if err != nil {
		t.Fatalf("InitServices() returned unexpected error: %v", err)
	}
}

func TestInitServices_NonDBStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	// Force config to non-db
	originalStorageType := config.AppConfig.StorageType
	config.AppConfig.StorageType = "file"
	defer func() { config.AppConfig.StorageType = originalStorageType }()

	// Call function under test
	log := logger.NewFakeLogger()
	err := initializer.InitServices(mockDB, log)
	if err != nil {
		t.Fatalf("InitServices() returned unexpected error: %v", err)
	}
	// No need to mock InitDBRepos because it should not be called
}
