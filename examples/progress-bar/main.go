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
	awf = alfred.NewWorkflow(
		alfred.WithOutStream(os.Stdout),
		alfred.WithLogStream(os.Stderr),
	)
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
	if awf.Cache(key).LoadItems(60*time.Second).Err() == nil {
		awf.Output()
		return
	}

	job := awf.Job(jobName)
	if !job.IsJob() && job.IsRunning() {
		awf.Rerun(0.5).Append(
			alfred.NewItem().Title("a background job is running"),
		)
		awf.Output()
		return
	}

	awf.Append(
		alfred.NewItem().Title("start a backgroup job"),
	).Rerun(0.5).Job(jobName).StartWithExit(os.Args[0], os.Args[1:]...)
	// clear existing(above) items as here is running as daemon
	awf.Clear()
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		awf.Append(
			alfred.NewItem().Title(fmt.Sprintf("%d", i)),
		)
	}
	awf.Cache(key).StoreItems().Workflow().Output()
	return
}
