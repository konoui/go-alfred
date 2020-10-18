package alfred

import (
	"errors"
	"os"
	"path/filepath"
)

// GetDataDir returns alfred data directory if data dir does not exist, creates it
func GetDataDir() (string, error) {
	dir := os.Getenv("alfred_workflow_data")
	if dir == "" {
		return "", errors.New("alfred_workflow_data is empty")
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(abs); err != nil {
		if err := os.MkdirAll(abs, os.ModePerm); err != nil {
			return "", err
		}
	}
	return abs, nil
}
