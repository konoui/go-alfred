package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/konoui/go-alfred"
)

var (
	awf *alfred.Workflow
)

const (
	cacheDir    = "./"
	dataDir     = "./"
	cacheSuffix = "-alfred-progress-bar.cache"
)

func main() {
	awf = alfred.NewWorkflow(
		alfred.WithLogLevel(alfred.LogLevelDebug),
		alfred.WithLogWriter(os.Stderr),
		alfred.WithCacheSuffix(cacheSuffix),
	)
	os.Exit(awf.Run(run))
}

func run(awf *alfred.Workflow) error {
	key := "test"
	jobName := "progress-bar"
	if awf.Cache(key).TTL(60*time.Second).LoadItems().Err() == nil {
		awf.Output()
		return nil
	}

	job := awf.Job(jobName)
	if !job.IsJob() && job.IsRunning() {
		awf.Rerun(0.5).Append(
			alfred.NewItem().Title("a background job is running"),
		)
		awf.Output()
		return nil
	}

	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	awf.Append(
		alfred.NewItem().Title("start a backgroup job"),
	).Rerun(0.5).Job(jobName).Logging().StartWithExit(cmd)
	// clear existing(above) items as here is running as daemon
	awf.Clear()
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		awf.Append(
			alfred.NewItem().Title(fmt.Sprintf("%d", i)),
		)
	}

	awf.Cache(key).StoreItems().Workflow().Output()
	return awf.Cache(key).Err()
}
