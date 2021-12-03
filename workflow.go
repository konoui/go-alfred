package alfred

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/konoui/go-alfred/update"
)

// Workflow is map of ScriptFilters
type Workflow struct {
	std        *ScriptFilter
	warn       *ScriptFilter
	err        *ScriptFilter
	system     *ScriptFilter
	cache      sync.Map
	markers    markers
	streams    *streams
	logger     *logger
	updater    Updater
	actions    []Initializer
	customEnvs *customEnvs
}

type streams struct {
	out io.Writer
	log io.Writer
}

type markers struct {
	done bool
}

type logger struct {
	l      Logger
	system Logger
	level  LogLevel
	tag    string
}

type customEnvs struct {
	autoUpdateTimeout time.Duration
	maxResults        int
	cacheSuffix       string
}

type Option func(*Workflow)

// NewWorkflow has simple ScriptFilter api
func NewWorkflow(opts ...Option) *Workflow {
	wf := &Workflow{
		std:    NewScriptFilter(),
		warn:   NewScriptFilter(),
		err:    NewScriptFilter(),
		system: NewScriptFilter(),
		streams: &streams{
			out: os.Stdout,
			log: io.Discard,
		},
		logger: &logger{
			tag:    "App",
			level:  LogLevelInfo,
			l:      nil,
			system: nil,
		},
		actions: []Initializer{new(envs), new(assets)},
		customEnvs: &customEnvs{
			autoUpdateTimeout: 2 * 60 * time.Second,
			maxResults:        0,
			cacheSuffix:       "",
		},
	}

	for _, opt := range opts {
		opt(wf)
	}

	if IsDebugEnabled() {
		wf.logger.level = LogLevelDebug
	}

	// sync log stream to logger
	wf.logger.l = newLogger(
		wf.streams.log,
		wf.logger.level,
		wf.logger.tag)
	wf.logger.system = newLogger(
		wf.streams.log,
		wf.logger.level,
		"System")
	return wf
}

// WithMaxResults arranges number of item result listed by Script Filter
func WithMaxResults(n int) Option {
	return func(wf *Workflow) {
		if n < 0 {
			return
		}
		if n > 0 {
			wf.customEnvs.maxResults = n
		}
	}
}

// WithLogLevel sets log level
func WithLogLevel(l LogLevel) Option {
	return func(wf *Workflow) {
		wf.logger.level = l
	}
}

// WithLogTag changes tag of a log message
// Log formats are [LogLevel][Tag] Message
func WithLogTag(tag string) Option {
	return func(wf *Workflow) {
		wf.logger.tag = tag
	}
}

// WithGitHubUpdater is managed github updater. updater will check newer version per `interval`
func WithGitHubUpdater(owner, repo, currentVersion string, interval time.Duration) Option {
	return WithUpdater(
		update.NewGitHubSource(
			owner,
			repo,
			currentVersion,
			update.WithCheckInterval(interval),
		),
	)
}

// WithUpdater supports native updater satisfing UpdaterSource interface
func WithUpdater(source update.UpdaterSource) Option {
	return func(wf *Workflow) {
		wf.updater = &updater{
			source: source,
			wf:     wf,
		}
		// add auto updater initializer
		wf.actions = append(wf.actions, new(autoUpdater))
	}
}

// WithInitializer registers Initializer.
// Initializer will be involved when OnInitialize() is called
func WithInitializer(i Initializer) Option {
	return func(wf *Workflow) {
		wf.actions = append(wf.actions, i)
	}
}

func WithLogWriter(w io.Writer) Option {
	return func(wf *Workflow) {
		wf.streams.log = w
	}
}

func WithOutWriter(w io.Writer) Option {
	return func(wf *Workflow) {
		wf.streams.out = w
	}
}

// WithAutoUpdateTimeout configures auto-update timeout for auto update Initializer
func WithAutoUpdateTimeout(v time.Duration) Option {
	return func(wf *Workflow) {
		wf.customEnvs.autoUpdateTimeout = v
	}
}

// WithCacheSuffix configures custom cacche siffux. default value is empty.
func WithCacheSuffix(suffix string) Option {
	return func(wf *Workflow) {
		wf.customEnvs.cacheSuffix = suffix
	}
}

func (w *Workflow) OutWriter() io.Writer {
	return w.streams.out
}

func (w *Workflow) LogWriter() io.Writer {
	return w.streams.log
}

// Append new items to ScriptFilter
func (w *Workflow) Append(item ...*Item) *Workflow {
	w.std.Items(item...)
	return w
}

// Clear items of ScriptFilter
func (w *Workflow) Clear() *Workflow {
	w.std.Clear()
	w.warn.Clear()
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
	w.warn.Items(
		NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(false).
			Icon(w.Assets().IconAlertNote()),
	)
	return w
}

// SetSystenInfo is useful for showing system information like update recommendation
// workflow ignores system information when store/load caches
func (w *Workflow) SetSystemInfo(i *Item) *Workflow {
	if i == nil {
		return w
	}
	w.system.Items(i)
	return w
}

func (w *Workflow) Bytes() []byte {
	if !w.err.IsEmpty() {
		return w.err.Bytes()
	}

	savedStdItems := make(Items, len(w.std.items), cap(w.std.items))
	copy(savedStdItems, w.std.items)
	savedWarnItems := make(Items, len(w.std.items), cap(w.std.items))
	copy(savedWarnItems, w.warn.items)
	defer func() {
		w.std.items = savedStdItems
		w.warn.items = savedWarnItems
	}()

	if w.isLimited() {
		w.std.items = savedStdItems[:w.customEnvs.maxResults]
	}

	if !w.system.IsEmpty() {
		items := w.std.items
		w.std.Clear()
		w.std.Items(w.system.items...)
		w.std.Items(items...)
		items = w.warn.items
		w.warn.Clear()
		w.warn.Items(w.system.items...)
		w.warn.Items(items...)
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
	w.err.Items(
		NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(false).
			Icon(w.Assets().IconCaution()),
	)
	w.Output()
	osExit(1)
}

// Output to io stream
func (w *Workflow) Output() *Workflow {
	if w.markers.done {
		w.sLogger().Warnln(sentMessage)
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
	limit := w.customEnvs.maxResults
	if limit > 0 && len(w.std.items) > limit {
		return true
	}
	return false
}
