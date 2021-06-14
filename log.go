package alfred

import (
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/logutils"
)

type Logger interface {
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Warnln(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
}

type myLogger struct {
	l *log.Logger
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

func (w *Workflow) Logger() Logger {
	return w.logger
}

func newLogger(out io.Writer, level LogLevel) Logger {
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{
			logutils.LogLevel(LogLevelDebug),
			logutils.LogLevel(LogLevelInfo),
			logutils.LogLevel(LogLevelWarn),
		},
		MinLevel: logutils.LogLevel(
			validate(level),
		),
		Writer: out,
	}

	return &myLogger{
		l: log.New(filter, "", 0),
	}
}

func validate(level LogLevel) LogLevel {
	switch level {
	case LogLevelWarn:
	case LogLevelInfo:
	case LogLevelDebug:
		return level
	}
	return LogLevelInfo
}

func (l *myLogger) Infof(format string, v ...interface{}) {
	level := string(LogLevelInfo)
	l.l.Printf("["+level+"] "+format, v...)
}

func (l *myLogger) Infoln(v ...interface{}) {
	level := string(LogLevelInfo)
	l.l.Printf("[" + level + "] " + fmt.Sprintln(v...))
}

func (l *myLogger) Debugf(format string, v ...interface{}) {
	level := string(LogLevelDebug)
	l.l.Printf("["+level+"] "+format, v...)
}

func (l *myLogger) Debugln(v ...interface{}) {
	level := string(LogLevelDebug)
	l.l.Printf("[" + level + "] " + fmt.Sprintln(v...))
}

func (l *myLogger) Warnln(v ...interface{}) {
	level := string(LogLevelWarn)
	l.l.Printf("[" + level + "] " + fmt.Sprintln(v...))
}

func (l *myLogger) Errorf(format string, v ...interface{}) {
	level := string(LogLevelError)
	l.l.Printf("["+level+"] "+format, v...)
}

func (l *myLogger) Errorln(v ...interface{}) {
	level := string(LogLevelError)
	l.l.Printf("[" + level + "] " + fmt.Sprintln(v...))
}
