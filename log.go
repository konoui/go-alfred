package alfred

import (
	"fmt"
	"io"
	"log"
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
	l     *log.Logger
	level LogLevel
	tag   string
}

type LogLevel int

const (
	LogLevelError LogLevel = iota + 1
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "Error"
	case LogLevelWarn:
		return "Warn"
	case LogLevelInfo:
		return "Info"
	case LogLevelDebug:
		return "Debug"
	default:
		return "Unknown"
	}
}

func (w *Workflow) sLogger() Logger {
	return w.logger.system
}

func (w *Workflow) Logger() Logger {
	return w.logger.l
}

func newLogger(out io.Writer, level LogLevel, tag string) Logger {
	return &myLogger{
		l:     log.New(out, "", 0),
		level: level,
		tag:   tag,
	}
}

func (l *myLogger) logInternal(logLevel LogLevel, tag, message string) {
	if l.level < logLevel {
		return
	}
	l.l.Printf("[%s][%s] %s", logLevel.String(), tag, message)
}

func (l *myLogger) Infof(format string, v ...interface{}) {
	l.logInternal(LogLevelInfo, l.tag, fmt.Sprintf(format, v...))
}

func (l *myLogger) Infoln(v ...interface{}) {
	l.Infof(fmt.Sprintln(v...))
}

func (l *myLogger) Debugf(format string, v ...interface{}) {
	l.logInternal(LogLevelDebug, l.tag, fmt.Sprintf(format, v...))
}

func (l *myLogger) Debugln(v ...interface{}) {
	l.Debugf(fmt.Sprintln(v...))
}

func (l *myLogger) Warnln(v ...interface{}) {
	l.logInternal(LogLevelWarn, l.tag, fmt.Sprintln(v...))
}

func (l *myLogger) Errorf(format string, v ...interface{}) {
	l.logInternal(LogLevelError, l.tag, fmt.Sprintf(format, v...))
}

func (l *myLogger) Errorln(v ...interface{}) {
	l.Errorf(fmt.Sprintln(v...))
}
