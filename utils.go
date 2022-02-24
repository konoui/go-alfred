package alfred

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"
)

const (
	// see https://www.alfredapp.com/help/workflows/script-environment-variables/
	envWorkflowData        = "alfred_workflow_data"
	envWorkflowCache       = "alfred_workflow_cache"
	envWorkflowBundleID    = "alfred_workflow_bundleid"
	envWorkflowDebug       = "alfred_debug"
	envWorkflowPreferences = "alfred_preferences"
	envWorkflowUID         = "alfred_workflow_uid"
)

var (
	// wrapper for tests
	osExit = os.Exit
	tmpDir = os.TempDir()
)

// GetBundleID returns value of alfred_workflow_bundleid environment variable
func (w *Workflow) GetBundleID() string {
	return os.Getenv(envWorkflowBundleID)
}

// GetDataDir returns value of alfred_workflow_data environment variable
func (w *Workflow) GetDataDir() string {
	return os.Getenv(envWorkflowData)
}

// GetCacheDir returns value of alfred_workflow_cache environment variable
func (w *Workflow) GetCacheDir() string {
	return os.Getenv(envWorkflowCache)
}

// GetWorkflowDir returns absolute path of the alfred workflow
func (w *Workflow) GetWorkflowDir() (string, error) {
	baseDir := os.Getenv(envWorkflowPreferences)
	if baseDir == "" {
		return "", fmt.Errorf(emptyEnvFormat, envWorkflowPreferences)
	}
	uid := os.Getenv(envWorkflowUID)
	if uid == "" {
		return "", fmt.Errorf(emptyEnvFormat, envWorkflowUID)
	}

	abs := filepath.Join(baseDir, "workflows", uid)
	if !PathExists(abs) {
		return "", fmt.Errorf("%s does not stat", abs)
	}
	return abs, nil
}

// IsDebugEnabled return true if alfred_debug is true
func IsDebugEnabled() bool {
	isDebug := parseBool(
		os.Getenv(envWorkflowDebug),
	)
	// debug env is highest priority
	return isDebug
}

// Normalize returns NFC string
// alfred workflow pass query as NFD
func Normalize(s string) string {
	return norm.NFC.String(s)
}

func parseBool(v string) bool {
	if strings.HasPrefix(v, "enable") {
		return true
	}
	if strings.HasPrefix(v, "disable") {
		return false
	}

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

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
