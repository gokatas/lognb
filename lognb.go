// Package lognb uses channels to implement non-blocking
// logging. Adapted from https://youtu.be/zDCKZn4-dck.
package lognb

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Logger struct {
	logs chan string
	wg   sync.WaitGroup
}

// New creates a logger that will write logs to w. Buf is the size of logs buffer.
func New(w io.Writer, buf int) *Logger {
	// New is sometimes called a factory function. It's useful
	// when you need to initialize one or more fields of a type.
	l := Logger{
		logs: make(chan string, buf),
	}

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for log := range l.logs {
			fmt.Fprintln(w, log)
		}
	}()

	return &l
}

func (l *Logger) Stop() {
	close(l.logs)
	l.wg.Wait()
}

// Write sends log to a channel to be written to logger's w if possible.
// Otherwise it writes a warning to stderr but doesn't block.
func (l *Logger) Write(log string) {
	select {
	case l.logs <- log:
	default:
		fmt.Fprintln(os.Stderr, "WARN: dropping logs on the floor")
	}
}
