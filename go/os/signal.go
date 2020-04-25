package main

// this test program is to test signal handler in golang and gracefully shutdown the program

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func shutdown() {
	fmt.Fprintf(os.Stderr, "%d exiting ...\n", os.Getpid())
	os.Exit(0)
}

func main() {
	// Prepare for graceful shutdown from os signals
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		for {
			select {
			case sig := <-osSignal:
				if sig == os.Kill || sig == syscall.SIGTERM {
					fmt.Fprintf(os.Stderr, "Got %s signal. Gracefully shutting down...\n", sig)
					shutdown()
				}
				if sig == os.Interrupt {
					fmt.Fprintf(os.Stderr, "Got %s signal. Dumping state to DB...\n", sig)
					shutdown()
				}
				if sig == syscall.SIGUSR1 {
					core := make([]int, 0)
					fmt.Printf("coredump: %s", core[0])
					// won't reach here
					shutdown()
				}
			}
		}
	}()

	t := 0
	for {
		fmt.Println("sleeping ..., waiting for signal", t)
		time.Sleep(5 * time.Second)
		t++
	}
}
