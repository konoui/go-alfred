package alfred

import (
	"path/filepath"

	"github.com/konoui/go-alfred/icon"
)

type Asseter interface {
	IconTrash() *Icon
	IconAlertNote() *Icon
	IconCaution() *Icon
	IconAlertStop() *Icon
	IconExec() *Icon
	Icon(string) *Icon
}

type Assets struct {
	_ *Workflow
}

func (w *Workflow) Asseter() Asseter {
	return w.assets
}

func getIconPath(filename string) string {
	return filepath.Join(icon.SystemIconPath, filename)
}

func (a *Assets) IconTrash() *Icon {
	return NewIcon().
		Path(getIconPath(icon.IconTrash))
}

func (a *Assets) IconAlertNote() *Icon {
	return NewIcon().
		Path(getIconPath(icon.IconAlerNote))
}

func (a *Assets) IconCaution() *Icon {
	return NewIcon().
		Path(getIconPath(icon.IconCaution))
}

func (a *Assets) IconAlertStop() *Icon {
	return NewIcon().
		Path(getIconPath(icon.IconAlertStop))
}

func (a *Assets) IconExec() *Icon {
	return NewIcon().
		Path(getIconPath(icon.IconExec))
}

func (a *Assets) Icon(filename string) *Icon {
	return NewIcon().
		Path(getIconPath(filename))
}
