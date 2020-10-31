package alfred

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	iconPath      = "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources"
	IconTrash     = NewIcon().SetPath(filepath.Join(iconPath, "TrashIcon.icns"))
	IconAlertNote = NewIcon().SetPath(filepath.Join(iconPath, "AlertNoteIcon.icns"))
	IconCaution   = NewIcon().SetPath(filepath.Join(iconPath, "AlertCautionIcon.icns"))
	IconAlertStop = NewIcon().SetPath(filepath.Join(iconPath, "AlertStopIcon.icns"))
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