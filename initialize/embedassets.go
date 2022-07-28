package initialize

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/konoui/go-alfred"
	"github.com/konoui/go-alfred/icon"
	"golang.org/x/sync/errgroup"
)

const (
	assetsDirName = "assets"
)

//go:embed assets/*
var embedAssetsFS embed.FS

type embedAssets struct {
	fallback alfred.Asseter
	wf       *alfred.Workflow
	dir      string
}

func NewEmbedAssets() alfred.Initializer {
	return &embedAssets{}
}

func GetAssetsDir(w *alfred.Workflow) string {
	return filepath.Join(w.GetDataDir(), assetsDirName)
}

// Condition returns true
// This means that the initializer is always executed
func (*embedAssets) Condition(*alfred.Workflow) bool { return true }

// Initialize creates asset files and directories
func (ea *embedAssets) Initialize(w *alfred.Workflow) (err error) {
	ea.fallback, ea.wf, ea.dir = w.Asseter(), w, GetAssetsDir(w)

	w.UpdateOpts(alfred.WithAsseter(ea))
	err = os.MkdirAll(ea.dir, os.ModePerm)
	if err != nil {
		return err
	}
	return generateAssets(ea.dir)
}

func (ea *embedAssets) getIcon(filename string, fallback *alfred.Icon) *alfred.Icon {
	path := filepath.Join(ea.dir, filename)
	if alfred.PathExists(path) {
		return alfred.NewIcon().
			Path(path)
	}
	return fallback
}

func (ea *embedAssets) IconTrash() *alfred.Icon {
	return ea.getIcon(icon.IconTrash, ea.fallback.IconTrash())
}

func (ea *embedAssets) IconAlertNote() *alfred.Icon {
	return ea.getIcon(icon.IconAlerNote, ea.fallback.IconAlertNote())
}

func (ea *embedAssets) IconCaution() *alfred.Icon {
	return ea.getIcon(icon.IconCaution, ea.fallback.IconCaution())
}

func (ea *embedAssets) IconAlertStop() *alfred.Icon {
	return ea.getIcon(icon.IconAlertStop, ea.fallback.IconAlertStop())
}

func (ea *embedAssets) IconExec() *alfred.Icon {
	return ea.getIcon(icon.IconExec, ea.fallback.IconExec())
}

func generateAssets(assetsDir string) error {
	icons, err := fs.Glob(embedAssetsFS, "**/*.icns")
	if err != nil {
		return err
	}

	var eg errgroup.Group
	for _, iconPath := range icons {
		relaPath := iconPath
		// Note relaPath format is `assets/<filename>`.
		// the assets is a dir name of go-alfred package, not `assetsDir` val.
		// remove directory name.
		name := filepath.Base(relaPath)
		path := filepath.Join(assetsDir, name)
		if alfred.PathExists(path) {
			continue
		}

		eg.Go(func() error {
			src, err := embedAssetsFS.Open(relaPath)
			if err != nil {
				return err
			}
			defer src.Close()

			dst, err := os.Create(path)
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				return err
			}
			return nil
		})
	}
	return eg.Wait()
}
