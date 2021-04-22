package main

import (
	"context"
	"os"

	"github.com/konoui/go-alfred"
)

func main() {
	awf := alfred.NewWorkflow(
		alfred.WithGitHubUpdater(
			"konoui", "alfred-tldr",
			"v0.0.1",
			0,
		),
		alfred.WithLogLevel(alfred.LogLevelDebug),
	)
	awf.SetLog(os.Stderr)

	if len(os.Args) >= 2 {
		if os.Args[1] == "--update" {
			_ = awf.Updater().Update(context.Background())
		}
	}

	if awf.Updater().NewerVersionAvailable(context.Background()) {
		awf.Append(
			alfred.NewItem().Title("newer version available!"),
		)
	}

	awf.Append(
		alfred.NewItem().Title("test"),
	)

	awf.Output()
}
