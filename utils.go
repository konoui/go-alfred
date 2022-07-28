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
