package logger

import (
	"io"
	"log"
)

func New(stream io.Writer) *log.Logger {
	return log.New(stream, "[*]", log.LstdFlags|log.Lshortfile)
}
