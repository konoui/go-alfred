package alfred

import (
	"fmt"
	"os"
)

// Initializer will invoke Initialize() when Condition returns true
type Initializer interface {
	Initialize(*Workflow) error
	Condition(*Workflow) bool
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

	w.actions = append(w.actions, initializers...)
	for _, i := range w.actions {
		if i == nil {
			continue
		}
		if i.Condition(w) {
			if err := i.Initialize(w); err != nil {
				return err
			}
		}
	}

	return nil
}

type normalizer struct{}

func (*normalizer) Condition(_ *Workflow) bool { return true }
func (*normalizer) Initialize(w *Workflow) (err error) {
	for idx, arg := range w.args {
		w.args[idx] = Normalize(arg)
	}
	return nil
}

type envs struct{}

// Condition returns true
// This means that the initializer is always executed
func (*envs) Condition(_ *Workflow) bool { return true }

// Initialize validates alfred workflow environment variables and creates directories
func (*envs) Initialize(w *Workflow) (err error) {
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
