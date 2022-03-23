package alfred

import (
	"fmt"
	"os"
)

// Initializer will invoke Initialize() when Condition returns true
type Initializer interface {
	Initialize(*Workflow) error
	Condition() bool
}

const emptyEnvFormat = "%s env is empty"

// OnInitialize executes pre-defined and custom initializers
// When using Run or Runsimple, do not need to involke OnInitialize.
func (w *Workflow) OnInitialize(initializers ...Initializer) error {
	if w.markers.initDone {
		w.sLogger().Warnln("The workflow has already initialized")
		return nil
	}
	defer func() { w.markers.initDone = true }()

	actions := append(w.actions, initializers...)
	for _, i := range actions {
		if i.Condition() {
			if err := i.Initialize(w); err != nil {
				return err
			}
		}
	}

	return nil
}

type normalizer struct{}

func (*normalizer) Condition() bool { return true }
func (*normalizer) Initialize(w *Workflow) (err error) {
	for idx, arg := range os.Args {
		os.Args[idx] = Normalize(arg)
	}
	return nil
}

type envs struct{}

// Condition returns true
// This means that the initializer is always executed
func (*envs) Condition() bool { return true }

// Initialize validates alfred workflow environment variables and creates directories
func (*envs) Initialize(w *Workflow) (err error) {
	defer func() {
		if w.customEnvs.skipEnvVerify && err != nil {
			w.sLogger().Warnln("skip environment initialization error")
			err = nil
		}
	}()
	bundleID := w.GetBundleID()
	if bundleID == "" {
		return fmt.Errorf(emptyEnvFormat, envWorkflowBundleID)
	}

	if err := initEnvDir(envWorkflowData); err != nil {
		return err
	}

	if err := initEnvDir(envWorkflowCache); err != nil {
		return err
	}
	return nil
}

func initEnvDir(key string) error {
	dir := os.Getenv(key)
	if dir == "" {
		return fmt.Errorf(emptyEnvFormat, key)
	}

	if !PathExists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
