package main

import (
	"fmt"
	"os"
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

func init() {
	awf = alfred.NewWorkflow()
	awf.SetOut(os.Stdout)
	awf.SetErr(os.Stderr)
	awf.SetCacheSuffix(cacheSuffix)
	if err := awf.SetCacheDir(cacheDir); err != nil {
		panic(err)
	}
	if err := awf.SetJobDir(dataDir); err != nil {
		panic(err)
	}
}

func main() {
	key := "test"
	jobName := "progress-bar"
	if awf.Cache(key).MaxAge(60*time.Second).LoadItems().Err() == nil {
		awf.Output()
		return
	}

	job := awf.Job(jobName)
	if !job.IsJob() && job.IsRunning() {
		awf.SetRerun(0.5).Append(
			alfred.NewItem().SetTitle("a background job is running"),
		)
		awf.Output()
		return
	}

	awf.Append(
		alfred.NewItem().SetTitle("start a backgroup job"),
	).SetRerun(0.5).Job(jobName).StartWithExit(os.Args[0], os.Args[1:]...)
	// clear existing(above) items as here is running as daemon
	awf.Clear()
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		awf.Append(
			alfred.NewItem().SetTitle(fmt.Sprintf("%d", i)),
		)
	}
	awf.Cache(key).StoreItems().Workflow().Output()
	return
}
