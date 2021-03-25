package main

import (
	"context"
	"os"

	"github.com/konoui/go-alfred"
	"github.com/konoui/go-alfred/update"
)

func main() {
	const current = "v0.0.1"
	awf := alfred.NewWorkflow(
		alfred.WithGitHubUpdater(
			"konoui", "alfred-tldr",
			update.WithVFormat(),
			update.WithCheckInterval(0),
		),
		alfred.WithLogLevel(alfred.LogLevelDebug),
	)
	awf.SetLog(os.Stderr)

	if len(os.Args) >= 2 {
		if os.Args[1] == "--update" {
			_ = awf.Updater().IfNewerVersionAvailable(current).Update(context.Background())
		}
	}

	awf.Updater().IfNewerVersionAvailable(current).AppendItem(
		alfred.NewItem().Title("newer version available!"),
	)

	awf.Append(
		alfred.NewItem().Title("test"),
	)

	awf.Output()
}
