package alfred

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// Initializer will invoke Initialize() when os.Args has Keyword()
// if Keyword() returns empty, Initialize() will be invoked
type Initializer interface {
	Initialize(*Workflow) error
	Keyword() string
}

const emptyEnvFormat = "%s env is empty"

var (
	updateTimeout = 3 * 60 * time.Second
)

var osExecutable = os.Executable

// OnInitialize executes followings
// 1. normalize arguments
// 2. execute pre-defined and custom initializers
// Custom initializer will be passed from arguments of OnInitialize or WithInitializer
func (w *Workflow) OnInitialize(initializers ...Initializer) error {
	for idx, arg := range os.Args {
		os.Args[idx] = Normalize(arg)
	}

	actions := append(w.actions, initializers...)
	for _, i := range actions {
		// If Keyword() returns empty, always do Initialize()
		if key := i.Keyword(); key == "" || hasArg(key) {
			if err := i.Initialize(w); err != nil {
				return err
			}
		}
	}

	return nil
}

type autoUpdater struct{}

// Keyword returns auto-update arguments
// This means that if the argument is specified, execute the Initializer
func (*autoUpdater) Keyword() string {
	return ArgWorkflowUpdate
}

// Initialize executes auto-updater of the workflow
func (*autoUpdater) Initialize(w *Workflow) error {
	jobName := "workflow-managed-update"
	if w.Job(jobName).IsRunning() {
		w.sLogger().Infoln("workflow-managed-update is already running")
		return nil
	}

	w.sLogger().Infoln("updating workflow...")
	self, err := osExecutable()
	if err != nil {
		return err
	}

	cmd := exec.Command(self, os.Args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdin = nil
	o, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	e, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	j, err := w.Job(jobName).Start(cmd)
	if err != nil {
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), updateTimeout)
	defer cancel()
	if j == JobWorker {
		err = w.Updater().Update(c)
		if err != nil {
			w.sLogger().Errorln("failed to update due to %v", err)
		}
	}

	if j == JobStarter {
		scanner := bufio.NewScanner(io.MultiReader(o, e))
		for scanner.Scan() {
			out := scanner.Text()
			w.sLogger().Infoln("[background-updater]", out)
		}
		if err := cmd.Wait(); err != nil {
			w.sLogger().Errorf("background-updater job failed due to %v. command dumps: %s", err, cmd.String())
			return fmt.Errorf("background-updater job failed: %w", err)
		}
		return nil
	}
	return nil
}

type assets struct{}

// Keyword returns empty string
// This means that the initializer is always executed
func (*assets) Keyword() string {
	return ""
}

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
	w.sLogger().Debugf("pre-defined assets are generated in %s", w.getAssetsDir())
	return nil
}

type envs struct{}

// Keyword returns empty string
// This means that the initializer is always executed
func (*envs) Keyword() string {
	return ""
}

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
	dir := getDir(key)
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
