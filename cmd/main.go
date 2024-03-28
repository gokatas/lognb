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

	"logger"
)

type device struct {
	problem bool
}

func (d *device) Write(p []byte) (int, error) {
	for d.problem {
		time.Sleep(time.Second)
	}
	return fmt.Print(string(p))
}

func main() {
	var d device
	l := logger.New(&d, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			for {
				l.Write(fmt.Sprintf("log from gr #%d", id))
				doSomething()
			}
		}(i)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	for {
		<-sigs
		d.problem = !d.problem
	}
}

func doSomething() {
	time.Sleep(time.Second)
}
