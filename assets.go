package alfred

import (
	"embed"
	"io/fs"
	"path/filepath"
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
	icons, err := fs.Glob(embedAssets, "**/*.png")
	if err != nil {
		return err
	}

	for _, relaPath := range icons {
		// Note remove directory name
		name := filepath.Base(relaPath)
		path := filepath.Join(assetsDir, name)
		if pathExists(path) {
			continue
		}

		// Note relaPath format is `assets/<filename>`
		// assets is a dir name of go-alfred package, not `assetsDirName` val
		data, err := embedAssets.ReadFile(relaPath)
		if err != nil {
			return err
		}

		if err := createFile(path, data); err != nil {
			return err
		}
	}

	return nil
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
		Path(a.getIconPath("TrashIcon.png"))
}

func (a *Assets) IconAlertNote() *Icon {
	return NewIcon().
		Path(a.getIconPath("AlertNoteIcon.png"))
}

func (a *Assets) IconCaution() *Icon {
	return NewIcon().
		Path(a.getIconPath("AlertCautionBadgeIcon.png"))
}

func (a *Assets) IconAlertStop() *Icon {
	return NewIcon().
		Path(a.getIconPath("AlertStopIcon.png"))
}

func (a *Assets) IconExec() *Icon {
	return NewIcon().
		Path(a.getIconPath("ExecutableBinaryIcon.icns"))
}
