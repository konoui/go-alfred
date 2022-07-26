package alfred

import (
	"path/filepath"
)

const (
	SystemIconPath = "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources"
	IconTrash      = "TrashIcon.icns"
	IconAlerNote   = "AlertNoteIcon.icns"
	IconCaution    = "AlertCautionBadgeIcon.icns"
	IconAlertStop  = "AlertStopIcon.icns"
	IconExec       = "ExecutableBinaryIcon.icns"
)

type Asseter interface {
	IconTrash() *Icon
	IconAlertNote() *Icon
	IconCaution() *Icon
	IconAlertStop() *Icon
	IconExec() *Icon
}

type Assets struct {
	_ *Workflow
}

func (w *Workflow) Asseter() Asseter {
	return w.assets
}

func getIconPath(filename string) string {
	return filepath.Join(SystemIconPath, filename)
}

func (a *Assets) IconTrash() *Icon {
	return NewIcon().
		Path(getIconPath(IconTrash))
}

func (a *Assets) IconAlertNote() *Icon {
	return NewIcon().
		Path(getIconPath(IconAlerNote))
}

func (a *Assets) IconCaution() *Icon {
	return NewIcon().
		Path(getIconPath(IconCaution))
}

func (a *Assets) IconAlertStop() *Icon {
	return NewIcon().
		Path(getIconPath(IconAlertStop))
}

func (a *Assets) IconExec() *Icon {
	return NewIcon().
		Path(getIconPath(IconExec))
}
