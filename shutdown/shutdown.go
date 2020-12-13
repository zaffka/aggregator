package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	lockWaitGroup   = &sync.WaitGroup{}
	lockWaitGroupCh = make(chan struct{})
	osInterruptCh   = make(chan os.Signal, 1)
	allDoneCh       = make(chan struct{})
)

// Lock is a function to prevent the main func to exit from execution till the OS interruption signal is received.
// It uses root context cancellation function to shut the context tree with a timeout.
//
// Also, it waits for the waitgroup app has under the hood (look at the WaitMe\ImDone funcs).
func Lock(rootCancelFn context.CancelFunc, timeOut time.Duration) <-chan struct{} {
	signal.Notify(osInterruptCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-osInterruptCh
		rootCancelFn()

		go func() {
			lockWaitGroup.Wait()
			close(lockWaitGroupCh)
		}()

		select {
		case <-time.After(timeOut):
		case <-lockWaitGroupCh:
		}

		close(allDoneCh)
	}()

	return allDoneCh
}

// WaitMe adds 1 to the package's waitgroup.
func WaitMe() {
	lockWaitGroup.Add(1)
}

// ImDone decrements package's waitgroup for 1.
func ImDone() {
	lockWaitGroup.Done()
}
