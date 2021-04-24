package alfred

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/text/unicode/norm"
)

var (
	iconPath      = "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources"
	IconTrash     = NewIcon().Path(filepath.Join(iconPath, "TrashIcon.icns"))
	IconAlertNote = NewIcon().Path(filepath.Join(iconPath, "AlertNoteIcon.icns"))
	IconCaution   = NewIcon().Path(filepath.Join(iconPath, "AlertCautionBadgeIcon.icns"))
	IconAlertStop = NewIcon().Path(filepath.Join(iconPath, "AlertStopIcon.icns"))
	IconExec      = NewIcon().Path(filepath.Join(iconPath, "ExecutableBinaryIcon.icns"))
)

const (
	// see https://www.alfredapp.com/help/workflows/script-environment-variables/
	workflowDataEnvKey        = "alfred_workflow_data"
	workflowCacheEnvKey       = "alfred_workflow_cache"
	workflowBundleIDEnvKey    = "alfred_workflow_bundleid"
	workflowDebugEnvKey       = "alfred_debug"
	workflowPreferencesEnvKey = "alfred_preferences"
	workflowUIDEnvKey         = "alfred_workflow_uid"
)

func (w *Workflow) GetBundleID() string {
	return os.Getenv(workflowBundleIDEnvKey)
}

func (w *Workflow) GetDataDir() string {
	return w.getDir(workflowDataEnvKey)
}

func (w *Workflow) GetWorkflowDir() (string, error) {
	baseDir := w.getDir(workflowPreferencesEnvKey)
	if baseDir == "" {
		return "", fmt.Errorf(emptyEnvFormat, workflowPreferencesEnvKey)
	}
	uid := os.Getenv(workflowUIDEnvKey)
	if uid == "" {
		return "", fmt.Errorf(emptyEnvFormat, workflowUIDEnvKey)
	}

	abs := filepath.Join(baseDir, "workflows", uid)
	if _, err := os.Stat(abs); err != nil {
		return "", fmt.Errorf("%s does not stat: %w", abs, err)
	}
	return abs, nil
}

func IsDebugEnabled() bool {
	isDebug := parseBool(
		os.Getenv(workflowDebugEnvKey),
	)
	// debug env is highest priority
	return isDebug
}

func (w *Workflow) getDir(key string) string {
	return os.Getenv(key)
}

// Normalize return NFC string
// alfred workflow pass query as NFD
func Normalize(s string) string {
	return norm.NFC.String(s)
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
