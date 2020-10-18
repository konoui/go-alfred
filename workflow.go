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
	std     ScriptFilter
	warn    ScriptFilter
	err     ScriptFilter
	cache   caches
	streams streams
	done    bool
	logger  *log.Logger
	dirs    map[string]string
}

type streams struct {
	out io.Writer
	err io.Writer
}

// SetOut redirect stdout
func (w *Workflow) SetOut(out io.Writer) {
	w.streams.out = out
}

// SetErr redirect stderr for debug util of the library.
func (w *Workflow) SetErr(stderr io.Writer) {
	w.streams.err = stderr
	w.logger = logger.New(w.streams.err)
}

// NewWorkflow has simple ScriptFilter api
func NewWorkflow() *Workflow {
	wf := &Workflow{
		std:  NewScriptFilter(),
		warn: NewScriptFilter(),
		err:  NewScriptFilter(),
		streams: streams{
			out: os.Stdout,
			err: ioutil.Discard,
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

// Rerun set rerun variable
func (w *Workflow) Rerun(i Rerun) *Workflow {
	w.std.Rerun = i
	w.warn.Rerun = i
	return w
}

// Variables set variables
func (w *Workflow) Variables(v Variables) *Workflow {
	w.std.Variables = v
	return w
}

// EmptyWarning create a new Item to Marshal　when there are no standard items
func (w *Workflow) EmptyWarning(title, subtitle string) *Workflow {
	w.warn = NewScriptFilter()
	w.warn.Append(
		&Item{
			Title:    title,
			Subtitle: subtitle,
			Valid:    true,
		})
	return w
}

// error append a new Item to error ScriptFilter
func (w *Workflow) error(title, subtitle string) *Workflow {
	w.err = NewScriptFilter()
	w.err.Append(
		&Item{
			Title:    title,
			Subtitle: subtitle,
			Valid:    true,
		})
	return w
}

// Marshal WorkFlow results
func (w *Workflow) Marshal() []byte {
	if len(w.std.Items) == 0 {
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
