package alfred

import (
	"fmt"
	"runtime/debug"
)

// RunSimple manages workflow environments and runs fn.
// This is useful when *Workflow is defined as global variable and fn refers it.
func (w *Workflow) RunSimple(fn func() error, i ...Initializer) (exitCode int) {
	wrraped := func(*Workflow) error { return fn() }
	return w.run(wrraped, i...)
}

// Run manages workflow environments and runs fn.
// Run will pass initialized *Workflow to argument of fn.
// This is useful for robust application.
func (w *Workflow) Run(fn func(*Workflow) error, i ...Initializer) (exitCode int) {
	return w.run(fn, i...)
}

func (w *Workflow) run(fn func(*Workflow) error, i ...Initializer) (exitCode int) {
	exitCode = 1
	if err := w.OnInitialize(i...); err != nil {
		outputErr(w, err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			w.sLogger().Errorln("Fatal Error")
			w.sLogger().Errorf("dump:\n %s\n", r)
			w.sLogger().Errorf("dump:\n %s\n", debug.Stack())

			err, ok := r.(error)
			if ok {
				outputErr(w, err)
				return
			}

			outputErr(w, fmt.Errorf("%v", r))
			return
		}
	}()

	if err := fn(w); err != nil {
		outputErr(w, err)
		return
	}

	return 0
}

func outputErr(w *Workflow, err error) {
	if err == nil {
		return
	}

	if w.markers.outputDone {
		return
	}

	w.err.Items(
		NewItem().
			Title(err.Error()).
			Subtitle("Please check workflow debug log").
			Icon(w.Assets().IconCaution()),
	)
	w.Output()
}
