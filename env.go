package alfred

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/konoui/go-alfred/env"
)

// GetBundleID returns value of alfred_workflow_bundleid environment variable
func (w *Workflow) GetBundleID() string {
	return os.Getenv(env.KeyWorkflowBundleID)
}

// GetDataDir returns value of alfred_workflow_data environment variable
func (w *Workflow) GetDataDir() string {
	return os.Getenv(env.KeyWorkflowData)
}

// GetCacheDir returns value of alfred_workflow_cache environment variable
func (w *Workflow) GetCacheDir() string {
	return os.Getenv(env.KeyWorkflowCache)
}

// GetWorkflowDir returns absolute path of the alfred workflow
func (w *Workflow) GetWorkflowDir() (string, error) {
	baseDir := os.Getenv(env.KeyWorkflowPreferences)
	if baseDir == "" {
		return "", fmt.Errorf(emptyEnvFormat, env.KeyWorkflowPreferences)
	}
	uid := os.Getenv(env.KeyWorkflowUID)
	if uid == "" {
		return "", fmt.Errorf(emptyEnvFormat, env.KeyWorkflowUID)
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
		os.Getenv(env.KeyWorkflowDebug),
	)
	// debug env is highest priority
	return isDebug
}
