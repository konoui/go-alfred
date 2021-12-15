package initialize

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/konoui/go-alfred"
)

const (
	ArgWorkflowUpdate = "workflow:update"
)

var (
	osExit       = os.Exit
	osExecutable = os.Executable
)

type updateChecker struct {
	timeout time.Duration
}

// HasUpdateArg return true if `ArgWorkflowUpdate` variable is specified
func hasUpdateArg() bool {
	for _, arg := range os.Args {
		if ArgWorkflowUpdate == arg {
			return true
		}
	}
	return false
}

func NewAutoUpdateChecker(timeout time.Duration) alfred.Initializer {
	return &updateChecker{timeout: timeout}
}

func (*updateChecker) Condition() bool { return true }
func (i *updateChecker) Initialize(w *alfred.Workflow) error {
	ctx, cancel := context.WithTimeout(context.Background(), i.timeout)
	defer cancel()
	if !hasUpdateArg() && w.Updater().IsNewVersionAvailable(ctx) {
		w.SetSystemInfo(
			alfred.NewItem().
				Title("New version workflow is available!").
				Subtitle("↩ for update").
				Autocomplete(ArgWorkflowUpdate).
				Valid(false).
				Icon(w.Assets().IconAlertNote()),
		)
	}
	return nil
}

type autoUpdater struct {
	timeout time.Duration
}

func NewAutoUpdater(timeout time.Duration) alfred.Initializer {
	return &autoUpdater{timeout: timeout}
}

func (*autoUpdater) Condition() bool {
	return hasUpdateArg()
}

// Initialize executes auto-updater of the workflow
func (i *autoUpdater) Initialize(w *alfred.Workflow) error {
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

	c, cancel := context.WithTimeout(context.Background(), i.timeout)
	defer cancel()
	switch j {
	case alfred.JobWorker:
		err = w.Updater().Update(c)
		code := 0
		if err != nil {
			w.Logger().Errorln("failed to update due to %v", err)
			code = 1
		}

		// after updating, worker process will exit
		osExit(code)
		return nil
	case alfred.JobStarter:
		scanner := bufio.NewScanner(io.MultiReader(o, e))
		for scanner.Scan() {
			out := scanner.Text()
			w.Logger().Infoln("[background-updater]", out)
		}

		if err := cmd.Wait(); err != nil {
			w.Logger().Errorf("background-updater job failed due to %v. command dumps: %s", err, cmd.String())
			return fmt.Errorf("background-updater job failed: %w", err)
		}

		// after waiting for worker process, output success message and exit
		w.SetSystemInfo(
			alfred.NewItem().
				Title("Update successfully"),
		).Output()
		osExit(0)
		return nil
	default:
		return fmt.Errorf("unexpected job status %d", j)
	}
}
