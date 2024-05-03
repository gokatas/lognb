// Package lognb uses channels to implement non-blocking
// logging. Adapted from https://youtu.be/zDCKZn4-dck.
package lognb

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const timeLayout = "2006/01/02 15:04:05"

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
		for s := range l.logs {
			fmt.Fprintf(w, "%s: %s\n", time.Now().Format(timeLayout), s)
		}
	}()

	return &l
}

func (l *Logger) Stop() {
	close(l.logs)
	l.wg.Wait()
}

// Print writes log to logger's w if possible. Otherwise it writes warning to
// stderr but doesn't block.
func (l *Logger) Print(s string) {
	select {
	case l.logs <- s:
	default:
		fmt.Fprintf(os.Stderr, "%s WARN dropping logs\n", time.Now().Format(timeLayout))
	}
}
