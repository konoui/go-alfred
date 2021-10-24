package alfred

import (
	"context"
	"errors"

	"github.com/konoui/go-alfred/update"
)

type updater struct {
	wf     *Workflow
	source update.UpdaterSource
}

type Updater interface {
	Update(context.Context) error
	NewerVersionAvailable(context.Context) bool
}

func (w *Workflow) Updater() Updater {
	if w.updater == nil {
		return &nilUpdater{}
	}
	return w.updater
}

func (u *updater) NewerVersionAvailable(ctx context.Context) bool {
	ok, err := u.source.NewerVersionAvailable(ctx)
	if err != nil {
		u.wf.sLogger().Warnln("failed to check newer version due to", err)
		return false
	}
	if ok {
		u.wf.sLogger().Infoln("newer version available")
		return true
	}
	u.wf.sLogger().Debugln("no newer version exists")
	return false
}

func (u *updater) Update(ctx context.Context) error {
	return u.source.IfNewerVersionAvailable().Update(ctx)
}

type nilUpdater struct{}

func (u *nilUpdater) Update(ctx context.Context) error {
	return errors.New("no implemented")
}

func (u *nilUpdater) NewerVersionAvailable(ctx context.Context) bool {
	return false
}
