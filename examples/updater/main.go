package main

import (
	"os"
	"time"

	"github.com/konoui/go-alfred"
	"github.com/konoui/go-alfred/initialize"
)

func main() {
	awf := alfred.NewWorkflow(
		alfred.WithGitHubUpdater(
			"konoui", "alfred-tldr",
			"v0.0.1",
			0,
		),
		alfred.WithLogLevel(alfred.LogLevelDebug),
		alfred.WithLogWriter(os.Stderr),
		alfred.WithInitializers(
			initialize.NewAutoUpdateChecker(2*time.Second),
			initialize.NewAutoUpdater(3*time.Minute),
		),
	)
	os.Exit(awf.Run(run))
}

func run(awf *alfred.Workflow) error {
	awf.Append(
		alfred.NewItem().Title("test"),
	)

	awf.Output()
	return nil
}
