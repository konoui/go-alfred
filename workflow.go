package alfred

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/konoui/go-alfred/logger"
)

// Workflow is map of ScriptFilters
type Workflow struct {
	std     *ScriptFilter
	warn    *ScriptFilter
	err     *ScriptFilter
	cache   caches
	streams streams
	done    bool
	logger  *log.Logger
	dirs    map[string]string
}

type streams struct {
	out io.Writer
}

// SetOut redirect stdout
func (w *Workflow) SetOut(out io.Writer) {
	w.streams.out = out
}

// SetErr redirect stderr for debug util of the library.
func (w *Workflow) SetErr(stderr io.Writer) {
	w.logger = logger.New(stderr)
}

// NewWorkflow has simple ScriptFilter api
func NewWorkflow() *Workflow {
	wf := &Workflow{
		std:  NewScriptFilter(),
		warn: NewScriptFilter(),
		err:  NewScriptFilter(),
		streams: streams{
			out: os.Stdout,
		},
		logger: logger.New(ioutil.Discard),
		dirs:   make(map[string]string),
	}

	return wf
}

// Append new items to ScriptFilter
func (w *Workflow) Append(item ...*Item) *Workflow {
	w.std.Append(item...)
	return w
}

// Clear items of ScriptFilter
func (w *Workflow) Clear() *Workflow {
	w.std.Clear()
	return w
}

// AddRerun add rerun variable
func (w *Workflow) Rerun(i Rerun) *Workflow {
	w.std.Rerun(i)
	w.warn.Rerun(i)
	w.err.Rerun(i)
	return w
}

// Variables add variables for ScriptFilter
func (w *Workflow) Variables(v Variables) *Workflow {
	w.std.Variables(v)
	return w
}

// Variable add variable for ScriptFilter
func (w *Workflow) Variable(key, value string) *Workflow {
	w.std.Variable(key, value)
	return w
}

// SetEmptyWarning message which will be showed if items is empty
func (w *Workflow) SetEmptyWarning(title, subtitle string) *Workflow {
	w.warn.Clear()
	w.warn.Append(
		NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(true).
			Icon(IconAlertNote),
	)
	return w
}

func (w *Workflow) error(title, subtitle string) *Workflow {
	w.err.Append(
		NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(true).
			Icon(IconCaution),
	)
	return w
}

// Marshal WorkFlow results
func (w *Workflow) Marshal() []byte {
	if w.IsEmpty() {
		return w.warn.Marshal()
	}
	return w.std.Marshal()
}

// Fatal output error to io stream and call os.Exit(1)
func (w *Workflow) Fatal(title, subtitle string) {
	if w.done {
		w.logger.Println(sentMessage)
		return
	}
	res := w.error(title, subtitle).err.Marshal()
	fmt.Fprintln(w.streams.out, string(res))
	os.Exit(1)
}

// Output to io stream
func (w *Workflow) Output() *Workflow {
	if w.done {
		w.logger.Println(sentMessage)
		return w
	}
	defer w.markDone()
	res := w.Marshal()
	fmt.Fprintln(w.streams.out, string(res))
	return w
}

// Logf print messages to debug console in workflow
func (w *Workflow) Logf(format string, v ...interface{}) *Workflow {
	w.logger.Printf(format, v...)
	return w
}

// IsEmpty return true if the items is empty
func (w *Workflow) IsEmpty() bool {
	return w.std.IsEmpty()
}

func (w *Workflow) markDone() {
	w.done = true
}
