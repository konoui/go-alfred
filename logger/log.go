package logger

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
}

type myLogger struct {
	l *log.Logger
}

const (
	LevelDebug string = "DEBUG"
	LevelInfo  string = "INFO"
	LevelWarn  string = "WARN"
)

func New(stream io.Writer, level string) Logger {
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{
			logutils.LogLevel(LevelDebug),
			logutils.LogLevel(LevelInfo),
			logutils.LogLevel(LevelWarn),
		},
		MinLevel: logutils.LogLevel(validate(LevelInfo)),
		Writer:   stream,
	}

	return &myLogger{
		l: log.New(filter, "", 0),
	}
}

func validate(level string) string {
	switch level {
	case LevelWarn:
	case LevelInfo:
	case LevelDebug:
		return level
	}
	return LevelInfo
}

func (l *myLogger) Infof(format string, v ...interface{}) {
	l.l.Printf("["+LevelInfo+"] "+format, v)
}

func (l *myLogger) Infoln(v ...interface{}) {
	l.l.Println("[" + LevelInfo + "] " + fmt.Sprintln(v...))
}

func (l *myLogger) Debugf(format string, v ...interface{}) {
	l.l.Printf("["+LevelDebug+"] "+format, v)
}

func (l *myLogger) Debugln(v ...interface{}) {
	l.l.Printf("[" + LevelDebug + "] " + fmt.Sprintln(v...))
}

func (l *myLogger) Warnln(v ...interface{}) {
	l.l.Println("[" + LevelWarn + "] " + fmt.Sprintln(v...))
}
