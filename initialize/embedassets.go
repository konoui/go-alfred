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
var embedSystemAssetsFS embed.FS

type embedAssets struct {
	fallback alfred.Asseter
	wf       *alfred.Workflow
	dir      string
	customFS []embed.FS
}

func NewEmbedAssets(customFS ...embed.FS) alfred.Initializer {
	return &embedAssets{
		customFS: customFS,
	}
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
	if err := generate(embedSystemAssetsFS, "**/*.icns", ea.dir, true); err != nil {
		return err
	}

	for _, f := range ea.customFS {
		if err := generate(f, "**/*", ea.dir, false); err != nil {
			return err
		}
	}
	return
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

func (ea *embedAssets) Icon(filename string) *alfred.Icon {
	return ea.getIcon(filename, ea.fallback.Icon(filename))
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
