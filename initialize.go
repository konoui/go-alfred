package alfred

import (
	"fmt"
	"os"

	"github.com/konoui/go-alfred/env"
)

var defaultInitializers = []Initializer{new(envs)}

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

type envs struct{}

// Condition returns true
// This means that the initializer is always executed
func (*envs) Condition(_ *Workflow) bool { return true }

// Initialize validates alfred workflow environment variables and creates directories
func (*envs) Initialize(w *Workflow) (err error) {
	bundleID := w.GetBundleID()
	if bundleID == "" {
		return fmt.Errorf(emptyEnvFormat, env.KeyWorkflowBundleID)
	}

	if err := initEnvDir(env.KeyWorkflowData); err != nil {
		return err
	}

	if err := initEnvDir(env.KeyWorkflowCache); err != nil {
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
