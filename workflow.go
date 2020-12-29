package alfred

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var tmpDir = os.TempDir()

// Workflow is map of ScriptFilters
type Workflow struct {
	std        *ScriptFilter
	warn       *ScriptFilter
	err        *ScriptFilter
	cache      caches
	streams    streams
	done       bool
	logger     Logger
	dirs       map[string]string
	maxResults int
}

type streams struct {
	out io.Writer
}

type Option func(*Workflow)

// NewWorkflow has simple ScriptFilter api
func NewWorkflow(opts ...Option) *Workflow {
	wf := &Workflow{
		std:  NewScriptFilter(),
		warn: NewScriptFilter(),
		err:  NewScriptFilter(),
		streams: streams{
			out: os.Stdout,
		},
		logger:     newLogger(ioutil.Discard, LogLevelInfo),
		dirs:       make(map[string]string),
		maxResults: 0,
	}

	for _, opt := range opts {
		opt(wf)
	}

	return wf
}

func WithMaxResults(n int) Option {
	return func(wf *Workflow) {
		if n < 0 {
			return
		}
		if n > 0 {
			wf.maxResults = n
		}
	}
}

func (w *Workflow) SetOut(out io.Writer) {
	w.streams.out = out
}

func (w *Workflow) SetLog(out io.Writer) {
	level := w.logger.logLevel()
	if IsDebugEnabled() {
		level = LogLevelDebug
	}
	w.logger = newLogger(out, level)

}

func (w *Workflow) SetLogLevel(level LogLevel) {
	if IsDebugEnabled() {
		level = LogLevelDebug
	}
	w.logger = newLogger(w.logger.Writer(), level)
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

// Rerun add rerun variable
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
	if len(w.err.items) != 0 {
		return w.err.Marshal()
	}

	if w.IsEmpty() {
		return w.warn.Marshal()
	}

	if limit := w.maxResults; limit > 0 && len(w.std.items) > limit {
		tmp := w.std.items
		w.std.items = w.std.items[:limit]
		defer func() {
			w.std.items = tmp
		}()
	}
	return w.std.Marshal()
}

// Fatal output error to io stream and call os.Exit(1)
func (w *Workflow) Fatal(title, subtitle string) {
	if w.done {
		w.logger.Infoln(sentMessage)
		return
	}

	res := w.error(title, subtitle).Marshal()
	fmt.Fprintln(w.streams.out, string(res))
	os.Exit(1)
}

// Output to io stream
func (w *Workflow) Output() *Workflow {
	if w.done {
		w.logger.Infoln(sentMessage)
		return w
	}
	defer w.markDone()
	res := w.Marshal()
	fmt.Fprintln(w.streams.out, string(res))
	return w
}

// IsEmpty return true if the items is empty
func (w *Workflow) IsEmpty() bool {
	return w.std.IsEmpty()
}

func (w *Workflow) markDone() {
	w.done = true
}
