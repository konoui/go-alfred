package alfred

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

const (
	assetsDirName  = "assets"
	systemIconPath = "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources"
)

//go:embed assets/*
var embedAssets embed.FS

type Assets struct {
	wf *Workflow
}

func (w *Workflow) getAssetsDir() string {
	return filepath.Join(w.GetDataDir(), assetsDirName)
}

func (w *Workflow) Assets() *Assets {
	return &Assets{
		wf: w,
	}
}

func generateAssets(assetsDir string) error {
	icons, err := fs.Glob(embedAssets, "**/*.icns")
	if err != nil {
		return err
	}

	var eg errgroup.Group
	for _, iconPath := range icons {
		relaPath := iconPath
		// Note relaPath format is `assets/<filename>`.
		// assets is a dir name of go-alfred package, not `assetsDirName` val.
		// remove directory name.
		name := filepath.Base(relaPath)
		path := filepath.Join(assetsDir, name)
		if pathExists(path) {
			continue
		}

		eg.Go(func() error {
			src, err := embedAssets.Open(relaPath)
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

func (a *Assets) getIconPath(filename string) string {
	path := filepath.Join(a.wf.getAssetsDir(), filename)
	if !pathExists(path) {
		return filepath.Join(systemIconPath, filename)
	}
	return path
}

func (a *Assets) IconTrash() *Icon {
	return NewIcon().
		Path(a.getIconPath("TrashIcon.icns"))
}

func (a *Assets) IconAlertNote() *Icon {
	return NewIcon().
		Path(a.getIconPath("AlertNoteIcon.icns"))
}

func (a *Assets) IconCaution() *Icon {
	return NewIcon().
		Path(a.getIconPath("AlertCautionBadgeIcon.icns"))
}

func (a *Assets) IconAlertStop() *Icon {
	return NewIcon().
		Path(a.getIconPath("AlertStopIcon.icns"))
}

func (a *Assets) IconExec() *Icon {
	return NewIcon().
		Path(a.getIconPath("ExecutableBinaryIcon.icns"))
}
