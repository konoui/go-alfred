package alfred

import (
	"fmt"
	"os"
)

const emptyEnvFormat = "%s env is empty"

func (w *Workflow) initDir(key string) error {
	dir := os.Getenv(key)
	if dir == "" {
		return fmt.Errorf(emptyEnvFormat, key)
	}

	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (w *Workflow) mustInit() {
	bundleID := w.GetBundleID()
	if bundleID == "" {
		panic(fmt.Errorf(emptyEnvFormat, workflowBundleIDEnvKey))
	}

	if err := w.initDir(workflowDataEnvKey); err != nil {
		panic(err)
	}

	if err := w.initDir(workflowCacheEnvKey); err != nil {
		panic(err)
	}
}

func (w *Workflow) Run(f func(wf *Workflow) error) {
	w.markManagedRun()
	defer func() {
		if err := recover(); err != nil {
			w.Logger().Errorln(err)
			w.Fatal("FATAL ERROR", "workflow console have more information")
		}
	}()

	w.mustInit()
	for idx, arg := range os.Args {
		os.Args[idx] = Normalize(arg)
	}

	if err := f(w); err != nil {
		w.Logger().Errorln(err)
		w.Fatal("application error", err.Error())
	}
}

func (w *Workflow) markManagedRun() {
	w.markers.managedRun = true
}
