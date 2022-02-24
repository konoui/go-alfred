package alfred

import (
	"path/filepath"
)

const (
	assetsDirName  = "assets"
	systemIconPath = "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources"
)

type Assets struct {
	wf *Workflow
}

func (w *Workflow) GetAssetsDir() string {
	return filepath.Join(w.GetDataDir(), assetsDirName)
}

func (w *Workflow) Assets() *Assets {
	return &Assets{
		wf: w,
	}
}

func (a *Assets) getIconPath(filename string) string {
	path := filepath.Join(a.wf.GetAssetsDir(), filename)
	if !PathExists(path) {
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
