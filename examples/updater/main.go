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
	if err := awf.OnInitialize(); err != nil {
		awf.Fatal(err.Error(), err.Error())
	}
	run(awf)
}

func run(awf *alfred.Workflow) error {
	if awf.Updater().NewerVersionAvailable(context.TODO()) {
		awf.Append(
			alfred.NewItem().Title("update workflow").
				Autocomplete(alfred.ArgWorkflowUpdate),
		)
	}

	awf.Append(
		alfred.NewItem().Title("test"),
	)

	awf.Output()
	return nil
}
