package update

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

var tmpDir = os.TempDir()

type UpdaterSource interface {
	NewerVersionAvailable(string) (bool, error)
	Update(ctx context.Context) error
}

type UpdaterSourceOption interface {
	setCheckInterval(time.Duration)
	enableVFormat()
}

type Option func(UpdaterSourceOption)

func WithCheckInterval(interval time.Duration) Option {
	return func(u UpdaterSourceOption) {
		u.setCheckInterval(interval)
	}
}

func WithVFormat() Option {
	return func(u UpdaterSourceOption) {
		u.enableVFormat()
	}
}

// compareVersions return true if `v2Str` is greater than `v1Str`
func compareVersions(v2Str, v1Str string, vFormat bool) (bool, error) {
	if vFormat {
		v2Str = strings.TrimPrefix(v2Str, "v")
		v1Str = strings.TrimPrefix(v1Str, "v")
	}

	v1, err := version.NewVersion(v1Str)
	if err != nil {
		return false, err
	}

	v2, err := version.NewVersion(v2Str)
	if err != nil {
		return false, err
	}

	if v2.GreaterThan(v1) {
		return true, nil
	}
	return false, nil
}

func updateContext(ctx context.Context, url string) error {
	filename := filepath.Base(url)
	path := filepath.Join(tmpDir, filename)
	if err := donwloadContext(ctx, url, path); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "open", path)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func donwloadContext(ctx context.Context, url, path string) error {
	if err, _ := os.Stat(path); err == nil {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	return nil
}

func hasAlfredWorkflowExt(s string) bool {
	return filepath.Ext(s) == ".alfredworkflow"
}