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
	*ScriptFilter
	warn       Items
	err        Items
	system     Items
	markers    markers
	streams    *streams
	logger     *logger
	updater    Updater
	actions    []Initializer
	customEnvs *customEnvs
	args       []string
}

type streams struct {
	out io.Writer
	log io.Writer
}

type markers struct {
	outputDone bool
	initDone   bool
}

type logger struct {
	l      Logger
	system Logger
	level  LogLevel
	tag    string
}

type customEnvs struct {
	maxResults  int
	cacheSuffix string
}

// Option is type for workflow configurations
type Option func(*Workflow)

// NewWorkflow has simple ScriptFilter api
func NewWorkflow(opts ...Option) *Workflow {
	wf := &Workflow{
		ScriptFilter: NewScriptFilter(),
		warn:         Items{},
		err:          Items{},
		system:       Items{},
		streams: &streams{
			out: os.Stdout,
			log: os.Stderr,
		},
		logger: &logger{
			tag:    "App",
			level:  LogLevelInfo,
			l:      nil,
			system: nil,
		},
		actions: defaultInitializers,
		customEnvs: &customEnvs{
			maxResults:  0,
			cacheSuffix: "",
		},
		args: normalizeAll(os.Args[1:]),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(wf)
	}

	wf.syncLogger()
	return wf
}

// UpdateOpts apply options to existing Workflow
func (w *Workflow) UpdateOpts(opts ...Option) *Workflow {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(w)
	}
	return w
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
		wf.syncLogger()
	}
}

// WithLogTag changes tag of a log message
// Log formats are [LogLevel][Tag] Message
func WithLogTag(tag string) Option {
	return func(wf *Workflow) {
		wf.logger.tag = tag
		wf.syncLogger()
	}
}

// WithGitHubUpdater is managed github updater. updater will check new version per `interval`
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
	}
}

// WithInitializers registers Initializer.
// Initializer will be involved on managed-run such as *Workflow.Run()
func WithInitializers(i ...Initializer) Option {
	return func(wf *Workflow) {
		a := make([]Initializer, 0, len(i)+len(defaultInitializers))
		a = append(a, defaultInitializers...)
		a = append(a, i...)
		wf.actions = a
	}
}

// WithLogWriter sets log output. os.Stderr is default value.
func WithLogWriter(w io.Writer) Option {
	return func(wf *Workflow) {
		wf.streams.log = w
		wf.syncLogger()
	}
}

// WithOutWriter sets ScriptFilter output. os.Stdout is default value.
func WithOutWriter(w io.Writer) Option {
	return func(wf *Workflow) {
		wf.streams.out = w
	}
}

// WithCacheSuffix configures custom cacche siffux. default value is empty.
func WithCacheSuffix(suffix string) Option {
	return func(wf *Workflow) {
		wf.customEnvs.cacheSuffix = suffix
	}
}

// WithArguments configures input args. default values are os.Args[1:]
func WithArguments(args ...string) Option {
	return func(w *Workflow) {
		w.args = normalizeAll(args)
	}
}

// OutWriter returns output writer
func (w *Workflow) OutWriter() io.Writer {
	return w.streams.out
}

// LogWriter returns logger writer
func (w *Workflow) LogWriter() io.Writer {
	return w.streams.log
}

func (w *Workflow) Args() []string {
	return w.args
}

// Append new items to ScriptFilter
func (w *Workflow) Append(item ...*Item) *Workflow {
	w.Items(item...)
	return w
}

// Rerun sets Rerun value
func (w *Workflow) Rerun(i Rerun) *Workflow {
	w.ScriptFilter.Rerun(i)
	return w
}

// Variables sets Variables for ScriptFilter
func (w *Workflow) Variables(v Variables) *Workflow {
	w.ScriptFilter.Variables(v)
	return w
}

// Variable sets Key/Value variable for ScriptFilter
func (w *Workflow) Variable(k, v string) *Workflow {
	w.ScriptFilter.Variable(k, v)
	return w
}

// Clear items of ScriptFilters
// Set* is not clear
func (w *Workflow) Clear() *Workflow {
	w.ScriptFilter.Clear()
	w.err = Items{}
	return w
}

// SetEmptyWarning displays messages if items are empty
func (w *Workflow) SetEmptyWarning(title, subtitle string) *Workflow {
	w.warn = append(w.warn,
		NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(false).
			Icon(IconAlertNote()),
	)
	return w
}

// SetSystenInfo is useful for displaying system information like update recommendation
// workflow ignores system information when it store/loads caches
func (w *Workflow) SetSystemInfo(i *Item) *Workflow {
	if i == nil {
		return w
	}
	w.system = append(w.system, i)
	return w
}

func (w *Workflow) Bytes() []byte {
	savedStdItems := make(Items, len(w.items), cap(w.items))
	savedErrItems := make(Items, len(w.err), cap(w.err))
	copy(savedStdItems, w.items)
	copy(savedErrItems, w.err)
	defer func() {
		w.items = savedStdItems
		w.err = savedErrItems
	}()

	if len(w.err) > 0 {
		items := w.err
		w.Clear()
		w.Items(items...)
		return w.ScriptFilter.Bytes()
	}

	if w.isLimited() {
		w.items = savedStdItems[:w.customEnvs.maxResults]
	}

	if len(w.system) > 0 {
		if w.IsEmpty() {
			w.Clear()
			w.Items(w.system...)
			w.Items(w.warn...)
			return w.ScriptFilter.Bytes()
		}

		items := w.items
		w.Clear()
		w.Items(w.system...)
		w.Items(items...)
		return w.ScriptFilter.Bytes()
	}

	if w.IsEmpty() {
		w.Clear()
		w.Items(w.warn...)
		return w.ScriptFilter.Bytes()
	}

	return w.ScriptFilter.Bytes()
}

// String show workflow outputs as JSON
func (w *Workflow) String() string {
	return string(w.Bytes())
}

// Fatal outputs error to io stream and call os.Exit(1)
func (w *Workflow) Fatal(title, subtitle string) {
	w.err = append(w.err,
		NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(false).
			Icon(IconCaution()))
	w.Output()
	osExit(1)
}

// Output outputs JSON of ScriptFilter to io stream
func (w *Workflow) Output() *Workflow {
	if w.markers.outputDone {
		w.sLogger().Warnln(sentMessage)
		return w
	}
	defer w.markDone()
	fmt.Fprintln(w.streams.out, w.String())
	return w
}

func (w *Workflow) markDone() {
	w.markers.outputDone = true
}

func (w *Workflow) isLimited() bool {
	// if maxResults equal 0, this means unlimited
	limit := w.customEnvs.maxResults
	if limit > 0 && len(w.items) > limit {
		return true
	}
	return false
}

func (w *Workflow) syncLogger() {
	if IsDebugEnabled() {
		w.logger.level = LogLevelDebug
	}

	// sync log stream to logger
	w.logger.l = newLogger(
		w.streams.log,
		w.logger.level,
		w.logger.tag)
	w.logger.system = newLogger(
		w.streams.log,
		w.logger.level,
		"System")
}
