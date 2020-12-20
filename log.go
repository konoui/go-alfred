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
	Writer() io.Writer
	logLevel() LogLevel
}

type myLogger struct {
	l     *log.Logger
	level LogLevel
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
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
		l:     log.New(filter, "", 0),
		level: level,
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

func (l *myLogger) Writer() io.Writer {
	return l.l.Writer()
}

func (l *myLogger) logLevel() LogLevel {
	return l.level
}

func (l *myLogger) Infof(format string, v ...interface{}) {
	level := string(LogLevelInfo)
	l.l.Printf("["+level+"] "+format, v)
}

func (l *myLogger) Infoln(v ...interface{}) {
	level := string(LogLevelInfo)
	l.l.Println("[" + level + "] " + fmt.Sprintln(v...))
}

func (l *myLogger) Debugf(format string, v ...interface{}) {
	level := string(LogLevelDebug)
	l.l.Printf("["+level+"] "+format, v)
}

func (l *myLogger) Debugln(v ...interface{}) {
	level := string(LogLevelDebug)
	l.l.Printf("[" + level + "] " + fmt.Sprintln(v...))
}

func (l *myLogger) Warnln(v ...interface{}) {
	level := string(LogLevelWarn)
	l.l.Println("[" + level + "] " + fmt.Sprintln(v...))
}
