// Log handles logging gracefully; the goroutines that do some important task
// (like sleeping) will not block just because it's not possible to write logs.
//
// Start 10 goroutines each of which will be writing logs to a device. Simulate
// a device problem by pressing Ctrl-C. Press Ctrl-C again to "fix" the problem.
// Ctrl-\ will terminate the program (with a core dump).
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"lognb"
)

type device struct {
	problem bool
}

func (d *device) Write(p []byte) (int, error) {
	for d.problem {
		time.Sleep(time.Second)
	}
	return fmt.Fprint(os.Stdout, string(p))
}

func main() {
	var d device
	// l := log.New(&d, log.Default().Prefix(), log.Default().Flags())
	l := lognb.New(&d, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				l.Print(fmt.Sprintf("goroutine %d does something", i))
				doSomething()
			}
		}(i)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		<-c
		d.problem = !d.problem
	}
}

func doSomething() { time.Sleep(time.Second) }
