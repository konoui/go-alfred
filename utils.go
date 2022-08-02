package alfred

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"
)

var (
	// wrapper for tests
	osExit = os.Exit
	tmpDir = os.TempDir()
)

// UnsetVariable unsets variable with key for existing Workflow
func UnsetVariable(w *Workflow, key string) {
	delete(w.variables, key)
}

// GetItems returns all appended items
func GetItems(w *Workflow) Items {
	return w.items
}

// ResetItems resets all items including EmptyWarning(), SetSystemInfo()
func ResetItems(w *Workflow) {
	w.items = Items{}
	w.system = Items{}
	w.warn = Items{}
	w.err = Items{}
}

// ResetSystemInfo resets items by SetSystemInfo()
func ResetSystemInfo(w *Workflow) {
	w.system = Items{}
}

// ResetWarning resets items by SetEmptyWarning()
func ResetWarning(w *Workflow) {
	w.warn = Items{}
}

// PathExists return true if path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Normalize returns NFC string
// alfred workflow pass query as NFD
func Normalize(s string) string {
	return norm.NFC.String(s)
}

func normalizeAll(args []string) []string {
	normargs := make([]string, len(args))
	for idx, arg := range args {
		normargs[idx] = Normalize(arg)
	}
	return normargs
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
