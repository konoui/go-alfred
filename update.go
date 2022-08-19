package alfred

import (
	"context"
	"errors"

	"github.com/konoui/go-alfred/internal/update"
)

type updater struct {
	wf     *Workflow
	source update.UpdaterSource
}

type Updater interface {
	Update(context.Context) error
	IsNewVersionAvailable(context.Context) bool
}

func (w *Workflow) Updater() Updater {
	if w.updater == nil {
		return &nilUpdater{}
	}
	return w.updater
}

func (u *updater) IsNewVersionAvailable(ctx context.Context) bool {
	ok, err := u.source.IsNewVersionAvailable(ctx)
	if err != nil {
		u.wf.sLogger().Warnln("failed to check new version due to", err)
		return false
	}
	if ok {
		u.wf.sLogger().Infoln("new workflow version available")
		return true
	}
	u.wf.sLogger().Debugln("new workflow version does not exist")
	return false
}

func (u *updater) Update(ctx context.Context) error {
	return u.source.IfNewVersionAvailable().Update(ctx)
}

type nilUpdater struct{}

func (u *nilUpdater) Update(ctx context.Context) error {
	return errors.New("no implemented")
}

func (u *nilUpdater) IsNewVersionAvailable(ctx context.Context) bool {
	return false
}
