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

const emptyEnvFormat = "%s env is empty"

var (
	checkTimeout  = 3 * time.Second
	updateTimeout = 3 * 60 * time.Second
)

var osExecutable = os.Executable

func (w *Workflow) OnInitialize() error {
	for idx, arg := range os.Args {
		os.Args[idx] = Normalize(arg)
	}

	if err := w.init(); err != nil {
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()
	if HasUpdateArg() && w.Updater().NewerVersionAvailable(c) {
		jobName := "workflow-managed-update"
		if w.Job(jobName).IsRunning() {
			w.Logger().Infoln("workflow-managed-update is already running")
			return nil
		}

		w.Logger().Infoln("updating workflow...")
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
				w.Logger().Errorln("failed to update due to %v", err)
			}
		}

		if j == JobStarter {
			scanner := bufio.NewScanner(io.MultiReader(o, e))
			for scanner.Scan() {
				out := scanner.Text()
				w.Logger().Infoln("[background-updater]", out)
			}
			if err := cmd.Wait(); err != nil {
				w.Logger().Errorln("background-updater job failed due to", err, ".", "command dumps:", cmd.String())
				return fmt.Errorf("background-updater job failed: %w", err)
			}
			return nil
		}
	}
	return nil
}

func initDir(key string) error {
	dir := getDir(key)
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

func (w *Workflow) init() error {
	bundleID := w.GetBundleID()
	if bundleID == "" {
		return fmt.Errorf(emptyEnvFormat, envWorkflowBundleID)
	}

	if err := initDir(envWorkflowData); err != nil {
		return err
	}

	if err := initDir(envWorkflowCache); err != nil {
		return err
	}
	return nil
}
