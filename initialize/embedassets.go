package initialize

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/konoui/go-alfred"
	"golang.org/x/sync/errgroup"
)

//go:embed assets/*
var embedAssetsFS embed.FS

type embedAssets struct{}

func NewEmbedAssets() alfred.Initializer {
	return &embedAssets{}
}

// Condition returns true
// This means that the initializer is always executed
func (*embedAssets) Condition(*alfred.Workflow) bool { return true }

// Initialize creates asset files and directories
func (*embedAssets) Initialize(w *alfred.Workflow) (err error) {
	err = os.MkdirAll(w.GetAssetsDir(), os.ModePerm)
	if err != nil {
		return err
	}
	return generateAssets(w.GetAssetsDir())
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
