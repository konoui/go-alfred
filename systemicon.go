package alfred

import (
	"path/filepath"
)

const (
	SystemIconPath = "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources"
	iconTrash      = "TrashIcon.icns"
	iconAlerNote   = "AlertNoteIcon.icns"
	iconCaution    = "AlertCautionBadgeIcon.icns"
	iconAlertStop  = "AlertStopIcon.icns"
	iconExec       = "ExecutableBinaryIcon.icns"
)

func getIconPath(filename string) string {
	return filepath.Join(SystemIconPath, filename)
}

func NewSystemIcon(filename string) *Icon {
	return NewIcon().
		Path(getIconPath(filename))
}

var (
	IconTrash     = func() *Icon { return NewSystemIcon(iconTrash) }
	IconAlertNote = func() *Icon { return NewSystemIcon(iconAlerNote) }
	IconCaution   = func() *Icon { return NewSystemIcon(iconCaution) }
	IconAlertStop = func() *Icon { return NewSystemIcon(iconAlertStop) }
	IconExec      = func() *Icon { return NewSystemIcon(iconExec) }
)
