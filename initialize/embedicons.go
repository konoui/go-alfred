package initialize

import (
	"embed"
	"path/filepath"

	"github.com/konoui/go-alfred"
)

//go:embed icons/*
var embedSystemAssetsFS embed.FS

type embedIcon struct {
	EmbedAsset
	fallback map[name]*alfred.Icon
}

type name int

const (
	trash name = iota
	alertNote
	caution
	alertStop
	execp
)

func NewEmbedSystemIcons() alfred.Initializer {
	fallback := map[name]*alfred.Icon{}
	fallback[trash] = alfred.IconTrash()
	fallback[alertNote] = alfred.IconAlertNote()
	fallback[caution] = alfred.IconCaution()
	fallback[alertStop] = alfred.IconAlertStop()
	fallback[execp] = alfred.IconExec()
	return &embedIcon{
		EmbedAsset: EmbedAsset{
			customFS: []embed.FS{embedSystemAssetsFS},
			flat:     true,
		},
		fallback: fallback,
	}
}

func Condition(w *alfred.Workflow) bool {
	return true
}

func (e *embedIcon) Initialize(w *alfred.Workflow) (err error) {
	err = e.EmbedAsset.Initialize(w)
	if err != nil {
		return err
	}

	alfred.IconTrash = e.IconTrash
	alfred.IconAlertNote = e.IconAlertNote
	alfred.IconAlertStop = e.IconAlertStop
	alfred.IconCaution = e.IconCaution
	alfred.IconExec = e.IconExec
	return nil
}

func (e *embedIcon) Down() {
	alfred.IconTrash = func() *alfred.Icon { return e.fallback[trash] }
	alfred.IconAlertNote = func() *alfred.Icon { return e.fallback[alertNote] }
	alfred.IconAlertStop = func() *alfred.Icon { return e.fallback[alertStop] }
	alfred.IconCaution = func() *alfred.Icon { return e.fallback[caution] }
	alfred.IconExec = func() *alfred.Icon { return e.fallback[execp] }
}

func (e *embedIcon) IconTrash() *alfred.Icon {
	return e.getIcon("TrashIcon.icns", e.fallback[trash])
}

func (e *embedIcon) IconAlertNote() *alfred.Icon {
	return e.getIcon("AlertNoteIcon.icns", e.fallback[alertNote])
}

func (e *embedIcon) IconCaution() *alfred.Icon {
	return e.getIcon("AlertCautionBadgeIcon.icns", e.fallback[caution])
}

func (e *embedIcon) IconAlertStop() *alfred.Icon {
	return e.getIcon("AlertStopIcon.icns", e.fallback[alertStop])
}

func (e *embedIcon) IconExec() *alfred.Icon {
	return e.getIcon("ExecutableBinaryIcon.icns", e.fallback[execp])
}

func (e *embedIcon) getIcon(filename string, fallback *alfred.Icon) *alfred.Icon {
	path := filepath.Join(e.dir, filename)
	if alfred.PathExists(path) {
		return alfred.NewIcon().
			Path(path)
	}
	return fallback
}
