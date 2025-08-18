package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var mu sync.Mutex

func ReadJSONFile[T any](filePath string) ([]T, error) {
	mu.Lock()
	defer mu.Unlock()

	if err := ensureDirExists(filepath.Dir(filePath)); err != nil {
		return nil, err
	}

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return []T{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var items []T
	err = json.Unmarshal(data, &items)
	return items, err
}

func WriteJSONFile[T any](filePath string, items []T) error {
	mu.Lock()
	defer mu.Unlock()

	if err := ensureDirExists(filepath.Dir(filePath)); err != nil {
		return err
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func ensureDirExists(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
