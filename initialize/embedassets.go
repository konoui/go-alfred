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

const (
	assetsDirName = "assets"
)

type EmbedAsset struct {
	dir      string
	customFS []embed.FS
	flat     bool
}

func NewEmbedAssets(customFS ...embed.FS) alfred.Initializer {
	return &EmbedAsset{
		customFS: customFS,
	}
}

func GetAssetsDir(w *alfred.Workflow) string {
	return filepath.Join(alfred.GetDataDir(), assetsDirName)
}

// Condition returns true
// This means that the initializer is always executed
func (*EmbedAsset) Condition(*alfred.Workflow) bool { return true }

// Initialize creates asset files and directories
func (ea *EmbedAsset) Initialize(w *alfred.Workflow) (err error) {
	ea.dir = GetAssetsDir(w)
	for _, f := range ea.customFS {
		if err := generate(f, "**/*", ea.dir, ea.flat); err != nil {
			return err
		}
	}
	return
}

func generate(fsys embed.FS, pattern, todir string, flat bool) error {
	blobs, err := fs.Glob(fsys, pattern)
	if err != nil {
		return err
	}

	var eg errgroup.Group
	for _, blobPath := range blobs {
		relaPath := blobPath

		path := filepath.Join(todir, relaPath)
		if flat {
			// Note relaPath format is `assets/<filename>`.
			// the assets is a dir name of go-alfred package, not `assetsDir` val.
			// remove directory name.
			name := filepath.Base(relaPath)
			path = filepath.Join(todir, name)
		}

		if alfred.PathExists(path) {
			continue
		}

		eg.Go(func() error {
			src, err := fsys.Open(relaPath)
			if err != nil {
				return err
			}
			defer src.Close()

			err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return err
			}

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
