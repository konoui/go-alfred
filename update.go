package alfred

import (
	"context"
	"errors"

	"github.com/konoui/go-alfred/update"
)

type updater struct {
	wf             *Workflow
	source         update.UpdaterSource
	currentVersion string
}

type Updater interface {
	Update(context.Context) error
	NewerVersionAvailable() bool
}

func (w *Workflow) Updater() Updater {
	if w.updater == nil {
		return &nilUpdater{}
	}
	return w.updater
}

func (u *updater) NewerVersionAvailable() bool {
	ok, err := u.source.NewerVersionAvailable(u.currentVersion)
	if err != nil {
		u.wf.Logger().Warnln("failed to check newer version", err)
		return false
	}
	if ok {
		u.wf.Logger().Infoln("newer version available!")
		return true
	}
	u.wf.Logger().Infoln("no newer version exists")
	return false
}

func (u *updater) Update(ctx context.Context) error {
	if u.NewerVersionAvailable() {
		return u.source.IfNewerVersionAvailable(u.currentVersion).Update(ctx)
	}
	return nil
}

type nilUpdater struct{}

func (u *nilUpdater) Update(ctx context.Context) error {
	return errors.New("no implemented")
}

func (u *nilUpdater) NewerVersionAvailable() bool {
	return false
}
