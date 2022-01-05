package alfred

import (
	"fmt"
	"os"
)

// Initializer will invoke Initialize() when os.Args has Keyword()
// if Keyword() returns empty, Initialize() will be invoked
type Initializer interface {
	Initialize(*Workflow) error
	Condition() bool
}

const emptyEnvFormat = "%s env is empty"

// OnInitialize executes followings
// 1. normalize arguments
// 2. execute pre-defined and custom initializers
// When using Run or Runsimple, do not need to involke OnInitialize.
func (w *Workflow) OnInitialize(initializers ...Initializer) error {
	if w.markers.initDone {
		w.sLogger().Warnln("The workflow has already initialized")
		return nil
	}
	defer func() { w.markers.initDone = true }()

	for idx, arg := range os.Args {
		os.Args[idx] = Normalize(arg)
	}

	actions := append(w.actions, initializers...)
	for _, i := range actions {
		// If Keyword() returns empty, always do Initialize()
		if i.Condition() {
			if err := i.Initialize(w); err != nil {
				return err
			}
		}
	}

	return nil
}

type assets struct{}

// Condition returns empty string
// This means that the initializer is always executed
func (*assets) Condition() bool { return true }

// Initialize generates/creates asset files and directories
func (*assets) Initialize(w *Workflow) (err error) {
	err = os.MkdirAll(w.getAssetsDir(), os.ModePerm)
	if err != nil {
		return err
	}

	err = generateAssets(w.getAssetsDir())
	if err != nil {
		return err
	}
	return nil
}

type envs struct{}

// Condition returns empty string
// This means that the initializer is always executed
func (*envs) Condition() bool { return true }

// Initialize validates alfred workflow environment variables and creates directories
func (*envs) Initialize(w *Workflow) error {
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

	if !pathExists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
