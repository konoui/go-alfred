package alfred

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/konoui/go-alfred/update"
)

// Workflow is map of ScriptFilters
type Workflow struct {
	std        *ScriptFilter
	warn       *ScriptFilter
	err        *ScriptFilter
	system     *ScriptFilter
	cache      caches
	streams    streams
	markers    markers
	logger     Logger
	maxResults int
	loglevel   LogLevel
	updater    Updater
}

type streams struct {
	out io.Writer
}

type markers struct {
	done bool
}

type Option func(*Workflow)

// NewWorkflow has simple ScriptFilter api
func NewWorkflow(opts ...Option) *Workflow {
	wf := &Workflow{
		std:    NewScriptFilter(),
		warn:   NewScriptFilter(),
		err:    NewScriptFilter(),
		system: NewScriptFilter(),
		streams: streams{
			out: os.Stdout,
		},
		logger:     newLogger(os.Stderr, LogLevelInfo),
		maxResults: 0,
		loglevel:   "",
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

func WithLogLevel(l LogLevel) Option {
	return func(wf *Workflow) {
		wf.loglevel = l
	}
}

// WithGitHubUpdater is managed github updater. updater will check newer version per `interval`
func WithGitHubUpdater(owner, repo, currentVersion string, interval time.Duration) Option {
	return func(wf *Workflow) {
		wf.updater = &updater{
			source: update.NewGitHubSource(
				owner,
				repo,
				currentVersion,
				update.WithCheckInterval(interval),
			),
			wf: wf,
		}
	}
}

// WithUpdater supports native updater satisfing UpdaterSource interface
func WithUpdater(source update.UpdaterSource, currentVersion string) Option {
	return func(wf *Workflow) {
		wf.updater = &updater{
			source: source,
			wf:     wf,
		}
	}
}

func (w *Workflow) SetOut(out io.Writer) {
	w.streams.out = out
}

func (w *Workflow) SetLog(out io.Writer) {
	level := LogLevelInfo
	if IsDebugEnabled() {
		level = LogLevelDebug
	} else if w.loglevel != "" {
		level = w.loglevel
	}
	w.logger = newLogger(out, level)
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
	w.warn.Append(
		NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(false).
			Icon(IconAlertNote),
	)
	return w
}

// SetSystenInfo is useful for showing system information like update recommendation
// workflow ignores system information when store/load caches
// item of icon will be overwritten with system icon.
func (w *Workflow) SetSystemInfo(i *Item) *Workflow {
	if i == nil {
		return w
	}
	w.system.Append(i)
	return w
}

func (w *Workflow) Bytes() []byte {
	if !w.err.IsEmpty() {
		return w.err.Bytes()
	}

	savedStdItems := make(Items, len(w.std.items))
	copy(savedStdItems, w.std.items)
	savedWarnItems := make(Items, len(w.warn.items))
	copy(savedWarnItems, w.warn.items)
	defer func() {
		w.std.items = savedStdItems
		w.warn.items = savedWarnItems
	}()

	if w.isLimited() {
		w.std.items = savedStdItems[:w.maxResults]
	}

	if !w.system.IsEmpty() {
		items := w.std.items
		w.std.Clear()
		w.std.Append(w.system.items...)
		w.std.Append(items...)
		items = w.warn.items
		w.warn.Clear()
		w.warn.Append(w.system.items...)
		w.warn.Append(items...)
	}

	if w.IsEmpty() {
		return w.warn.Bytes()
	}

	return w.std.Bytes()
}

func (w *Workflow) String() string {
	return string(w.Bytes())
}

var osExit = os.Exit

// Fatal output error to io stream and call os.Exit(1)
func (w *Workflow) Fatal(title, subtitle string) {
	w.err.Append(
		NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(false).
			Icon(IconCaution),
	)
	w.Output()
	osExit(1)
}

// Output to io stream
func (w *Workflow) Output() *Workflow {
	if w.markers.done {
		w.Logger().Warnln(sentMessage)
		return w
	}
	defer w.markDone()
	fmt.Fprintln(w.streams.out, w.String())
	return w
}

// IsEmpty return true if the items is empty
func (w *Workflow) IsEmpty() bool {
	return w.std.IsEmpty()
}

func (w *Workflow) markDone() {
	w.markers.done = true
}

func (w *Workflow) isLimited() bool {
	// if maxResults equal 0, this means unlimited
	limit := w.maxResults
	if limit > 0 && len(w.std.items) > limit {
		return true
	}
	return false
}
