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
	// EnvAlfredAutoUpdateWorkflow is bool value
	EnvAutoUpdateWorkflow = "alfred_auto_update_workflow"
	ArgWorkflowUpdate     = "workflow:update"
)

func (w *Workflow) GetBundleID() string {
	return os.Getenv(envWorkflowBundleID)
}

func (w *Workflow) GetDataDir() string {
	return getDir(envWorkflowData)
}

func (w *Workflow) GetCacheDir() string {
	return getDir(envWorkflowCache)
}

func (w *Workflow) GetWorkflowDir() (string, error) {
	baseDir := getDir(envWorkflowPreferences)
	if baseDir == "" {
		return "", fmt.Errorf(emptyEnvFormat, envWorkflowPreferences)
	}
	uid := os.Getenv(envWorkflowUID)
	if uid == "" {
		return "", fmt.Errorf(emptyEnvFormat, envWorkflowUID)
	}

	abs := filepath.Join(baseDir, "workflows", uid)
	if !pathExists(abs) {
		return "", fmt.Errorf("%s does not stat", abs)
	}
	return abs, nil
}

// IsAutoUpdateWorkflowEnabled return false only when env value is false
// otherwise return true e.g.) env is not set
func IsAutoUpdateWorkflowEnabled() bool {
	v := os.Getenv(EnvAutoUpdateWorkflow)
	if v == "" {
		return true
	}
	return parseBool(v)
}

func HasUpdateArg() bool {
	return hasArg(ArgWorkflowUpdate)
}

func IsDebugEnabled() bool {
	isDebug := parseBool(
		os.Getenv(envWorkflowDebug),
	)
	// debug env is highest priority
	return isDebug
}

func getDir(key string) string {
	return os.Getenv(key)
}

// Normalize return NFC string
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

func hasArg(v string) bool {
	for _, arg := range os.Args {
		if arg == v {
			return true
		}
	}
	return false
}

func createFile(path string, data []byte) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return
	}
	return
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
