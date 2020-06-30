package logger

import (
	"io"
	"log"
)

func New(stream io.Writer) *log.Logger {
	return log.New(stream, "[INFO]", log.LstdFlags|log.Lshortfile)
}
