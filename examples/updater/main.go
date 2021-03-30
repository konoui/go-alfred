package main

import (
	"context"
	"os"

	"github.com/konoui/go-alfred"
	"github.com/konoui/go-alfred/update"
)

func main() {
	awf := alfred.NewWorkflow(
		alfred.WithGitHubUpdater(
			"konoui", "alfred-tldr",
			"v0.0.1",
			update.WithCheckInterval(0),
		),
		alfred.WithLogLevel(alfred.LogLevelDebug),
	)
	awf.SetLog(os.Stderr)

	if len(os.Args) >= 2 {
		if os.Args[1] == "--update" {
			_ = awf.Updater().Update(context.Background())
		}
	}

	if awf.Updater().NewerVersionAvailable() {
		awf.Append(
			alfred.NewItem().Title("newer version available!"),
		)
	}

	awf.Append(
		alfred.NewItem().Title("test"),
	)

	awf.Output()
}
