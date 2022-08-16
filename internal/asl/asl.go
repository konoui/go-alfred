package asl

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Config struct {
	timeout           time.Duration
	loggerCommandName string
}

type asl struct {
	bin string
	cfg *Config
	mux sync.Mutex
}

type Option func(*Config)

// WithCommandTimeout configures logger command timeout
func WithCommandTimeout(v time.Duration) Option {
	return func(a *Config) {
		a.timeout = v
	}
}

// WithCommandTimeout configures logger command name.
// default is `logger`. This is useful for test mock
func WithLoggerCommandName(s string) Option {
	return func(a *Config) {
		a.loggerCommandName = s
	}
}

// New returns io.Writer for apple system logger(syslogd)
// io.Writer.Write() writes data to asl by `logger` command
// `log show` command can dig logs.
// e.g.) `log show --style compact --info --predicate 'process == "logger"'`
func New(opts ...Option) (io.Writer, error) {
	cfg := &Config{
		timeout:           500 * time.Millisecond,
		loggerCommandName: "logger",
	}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(cfg)
	}

	bin, err := exec.LookPath(cfg.loggerCommandName)
	if err != nil {
		// default logger command path
		bin = "/usr/bin/" + cfg.loggerCommandName
		if _, err := os.Stat(bin); err != nil {
			return nil, err
		}
	}

	a := &asl{
		bin: bin,
		cfg: cfg,
	}

	return log.New(a, "", 0).Writer(), nil
}

func (a *asl) log(msg string) error {
	a.mux.Lock()
	defer a.mux.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, a.bin, msg) //nolint
	return cmd.Run()
}

func (a *asl) Write(b []byte) (n int, err error) {
	return len(b), a.log(string(b))
}
