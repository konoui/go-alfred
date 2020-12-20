package alfred

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

var (
	iconPath      = "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources"
	IconTrash     = NewIcon().Path(filepath.Join(iconPath, "TrashIcon.icns"))
	IconAlertNote = NewIcon().Path(filepath.Join(iconPath, "AlertNoteIcon.icns"))
	IconCaution   = NewIcon().Path(filepath.Join(iconPath, "AlertCautionIcon.icns"))
	IconAlertStop = NewIcon().Path(filepath.Join(iconPath, "AlertStopIcon.icns"))
)

// GetDataDir returns alfred data directory if data dir does not exist, creates it
func GetDataDir() (string, error) {
	dir := os.Getenv("alfred_workflow_data")
	if dir == "" {
		return "", errors.New("alfred_workflow_data env is empty")
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

func isEnabledDebug() bool {
	isDebug := parseBool(
		os.Getenv("alfred_debug"),
	)
	// debug env is highest priority
	return isDebug
}

func parseBool(v string) bool {
	i, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return i == 1
	}

	b, err := strconv.ParseBool(v)
	if err == nil {
		return b
	}

	return false
}
