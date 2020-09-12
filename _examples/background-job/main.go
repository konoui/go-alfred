package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/konoui/go-alfred"
)

var (
	awf *alfred.Workflow
)

const (
	dataDir = "./data"
)

func init() {
	awf = alfred.NewWorkflow()
	awf.SetErr(os.Stderr)
	if err := awf.SetJobDir(dataDir); err != nil {
		panic(err)
	}
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if strings.EqualFold(getQuery(os.Args, 1), "start") {
		return startJobs()
	}
	if strings.EqualFold(getQuery(os.Args, 1), "kill") {
		return terminateJob(getQuery(os.Args, 2))
	}
	return listJobs()
}

func getQuery(args []string, idx int) string {
	if len(args) > idx {
		return args[idx]
	}
	return ""
}

func startJobs() error {
	jobName := "backgound-job"
	awf.Job(jobName).Logging().StartWithExit(os.Args[0], os.Args[1:]...)
	// next instructions will be executed as job
	awf.Clear()
	return runCmd()
}

func runCmd() error {
	cmd := exec.Command("sleep", "300")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func listJobs() error {
	awf.EmptyWarning("no jobs", "")
	jobs := awf.ListJobs()
	for _, job := range jobs {
		awf.Append(
			alfred.NewItem().SetTitle(job.Name()),
		)
	}
	awf.Output()
	return nil
}

func terminateJob(jobName string) error {
	return awf.Job(jobName).Terminate()
}
