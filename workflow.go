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

// Append a new Item to standard ScriptFilter
func (w *Workflow) Append(item *Item) *Workflow {
	w.std.Append(item)
	return w
}

// Clear items of standard ScriptFilter
func (w *Workflow) Clear() *Workflow {
	w.std.Items = Items{}
	return w
}

// SetRerun set rerun variable
func (w *Workflow) SetRerun(i Rerun) *Workflow {
	w.std.SetRerun(i)
	w.warn.SetRerun(i)
	w.err.SetRerun(i)
	return w
}

// SetVariables set variables for ScriptFilter
func (w *Workflow) SetVariables(v Variables) *Workflow {
	w.std.SetVariables(v)
	return w
}

// SetVariable set variable for ScriptFilter
func (w *Workflow) SetVariable(key, value string) *Workflow {
	w.std.SetVariable(key, value)
	return w
}

// EmptyWarning create a new Item to Marshal　when there are no standard items
func (w *Workflow) EmptyWarning(title, subtitle string) *Workflow {
	w.warn = NewScriptFilter()
	w.warn.Append(
		NewItem().
			SetTitle(title).
			SetSubtitle(subtitle).
			SetValid(true).
			SetIcon(IconAlertNote),
	)
	return w
}

// error append a new Item to error ScriptFilter
func (w *Workflow) error(title, subtitle string) *Workflow {
	w.err = NewScriptFilter()
	w.err.Append(
		NewItem().
			SetTitle(title).
			SetSubtitle(subtitle).
			SetValid(true).
			SetIcon(IconCaution),
	)
	return w
}

// Marshal WorkFlow results
func (w *Workflow) Marshal() []byte {
	if w.std.IsEmpty() {
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
	w.error(title, subtitle)
	res := w.err.Marshal()
	fmt.Fprintln(w.streams.out, string(res))
	w.done = true
	os.Exit(1)
}

// Output to io stream
func (w *Workflow) Output() *Workflow {
	if w.done {
		w.logger.Println(sentMessage)
		return w
	}
	res := w.Marshal()
	fmt.Fprintln(w.streams.out, string(res))
	w.done = true
	return w
}

// Logf print messages to debug console in workflow
func (w *Workflow) Logf(format string, v ...interface{}) *Workflow {
	w.logger.Printf(format, v...)
	return w
}
