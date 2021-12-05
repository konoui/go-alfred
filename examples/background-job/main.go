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
	jobName = "backgound-job"
)

func init() {
	awf = alfred.NewWorkflow(
		alfred.WithLogLevel(alfred.LogLevelDebug),
		alfred.WithLogWriter(os.Stderr),
	)
}

func main() {
	if err := awf.OnInitialize(); err != nil {
		panic(err)
	}
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if strings.EqualFold(getQuery(os.Args, 1), "start") {
		return startJobs()
	}
	if strings.EqualFold(getQuery(os.Args, 1), "kill") {
		return terminateJob(jobName)
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
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	awf.Logger().Infof("starting the %s ...", jobName)
	awf.Job(jobName).StartWithExit(cmd)
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
	awf.SetEmptyWarning("no jobs", "")
	awf.Logger().Infoln("listing jobs ...")
	jobs := awf.ListJobs()
	for _, job := range jobs {
		awf.Append(
			alfred.NewItem().Title(job.Name()),
		)
	}
	awf.Output()
	return nil
}

func terminateJob(jobName string) error {
	awf.Logger().Infof("kill the %s ...", jobName)
	return awf.Job(jobName).Terminate()
}
