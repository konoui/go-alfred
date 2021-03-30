package alfred

import (
	"context"

	"github.com/konoui/go-alfred/update"
)

type updater struct {
	wf             *Workflow
	source         update.UpdaterSource
	currentVersion string
}

type UpdaterSource interface {
	IfNewerVersionAvailable(string) Updater
}

type Updater interface {
	Update(context.Context) error
	AppendItem(...*Item)
}

func (w *Workflow) Updater() UpdaterSource {
	return &updater{
		wf:     w,
		source: w.updater,
	}
}

func (u *updater) IfNewerVersionAvailable(currentVersion string) Updater {
	u.currentVersion = currentVersion
	return u
}

func (u *updater) Update(ctx context.Context) error {
	ok, err := u.source.NewerVersionAvailable(u.currentVersion)
	if err != nil {
		u.wf.Logger().Warnln("failed to check newer version", err)
		return err
	}
	if ok {
		return u.source.Update(ctx)
	}
	return nil
}

func (u *updater) AppendItem(items ...*Item) {
	ok, err := u.source.NewerVersionAvailable(u.currentVersion)
	if err != nil {
		u.wf.Logger().Warnln("failed to check newer version", err)
		return
	}
	if ok {
		u.wf.Append(items...)
		u.wf.logger.Infoln("newer version available!")
		return
	}
	u.wf.logger.Infoln("no newer version exists")
}
