package alfred

import (
	"fmt"
	"io"
)

var (
	_ Appender      = (*Workflow)(nil)
	_ Outputer      = (*Workflow)(nil)
	_ Clearer       = (*Workflow)(nil)
	_ Setter        = (*Workflow)(nil)
	_ IO            = (*Workflow)(nil)
	_ fmt.Stringer  = (*Workflow)(nil)
	_ Hooker        = (*Workflow)(nil)
	_ LogGetter     = (*Workflow)(nil)
	_ UpdateGetter  = (*Workflow)(nil)
	_ CacheGetter   = (*Workflow)(nil)
	_ AssetGetter   = (*Workflow)(nil)
	_ Runner        = (*Workflow)(nil)
	_ EnvGetter     = (*Workflow)(nil)
	_ ArgGetter     = (*Workflow)(nil)
	_ Filter        = (*Workflow)(nil)
	_ OptionUpdater = (*Workflow)(nil)
	_ JobGetter     = (*Workflow)(nil)
)

type Appender interface {
	Append(item ...*Item) *Workflow
	Rerun(i Rerun) *Workflow
	Variables(Variables) *Workflow
	Variable(key, value string) *Workflow
}

type Setter interface {
	SetEmptyWarning(title, subtitle string) *Workflow
	SetSystemInfo(i *Item) *Workflow
}

type Outputer interface {
	Output() *Workflow
	Fatal(title, subtitle string)
}

type Clearer interface {
	Clear() *Workflow
	IsEmpty() bool
}

type IO interface {
	OutWriter() io.Writer
	LogWriter() io.Writer
}

type Hooker interface {
	OnInitialize(initializers ...Initializer) error
}

type LogGetter interface {
	Logger() Logger
}

type UpdateGetter interface {
	Updater() Updater
}

type CacheGetter interface {
	Cache(key string) Cacher
}

type AssetGetter interface {
	Asseter() Asseter
}

type Runner interface {
	RunSimple(fn func() error, i ...Initializer) (exitCode int)
	Run(fn func(*Workflow) error, i ...Initializer) (exitCode int)
}

type EnvGetter interface {
	GetBundleID() string
	GetDataDir() string
	GetCacheDir() string
	GetWorkflowDir() (string, error)
	GetAssetsDir() string
}

type ArgGetter interface {
	Args() []string
}

type Filter interface {
	Filter(query string) *Workflow
	FilterByItemProperty(f func(s string) bool, property ItemProperty) *Workflow
}

type OptionUpdater interface {
	UpdateOpts(opts ...Option) *Workflow
}

// TODO implement
type JobGetter interface {
	Job(name string) *Job
	ListJobs() []*Job
}
