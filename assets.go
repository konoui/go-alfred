package alfred

import (
	"path/filepath"
)

const (
	assetsDirName  = "assets"
	systemIconPath = "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources"
)

type Asseter interface {
	IconTrash() *Icon
	IconAlertNote() *Icon
	IconCaution() *Icon
	IconAlertStop() *Icon
	IconExec() *Icon
}

type Assets struct {
	wf *Workflow
}

func (w *Workflow) GetAssetsDir() string {
	return filepath.Join(w.GetDataDir(), assetsDirName)
}

func (w *Workflow) Asseter() Asseter {
	return &Assets{
		wf: w,
	}
}

func (a *Assets) getIconPath(filename string) string {
	path := filepath.Join(a.wf.GetAssetsDir(), filename)
	if PathExists(path) {
		return path
	}
	return filepath.Join(systemIconPath, filename)
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
