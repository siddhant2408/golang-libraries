// Package ctxsignal provides a bridge for context and POSIX signals.
package ctxsignal

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/siddhant2408/golang-libraries/goroutine"
)

// RegisterCancel registers a listener for the SIGINT/SIGTERM signals.
// The first time that the signal is received, the context is canceled.
// The second time, the program exits with status 0.
func RegisterCancel(parent context.Context) (ctx context.Context, unregister func()) {
	ctx, cancel := context.WithCancel(parent)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	waitWatchSignal := goroutine.Go(func() {
		watchSignal(c, cancel)
	})
	unregister = newUnregisterFunc(c, waitWatchSignal)
	return ctx, unregister
}

func watchSignal(c <-chan os.Signal, cancel func()) {
	sig, ok := <-c
	if !ok {
		return
	}
	log.Printf("%q signal received, context canceled. Send it again to exit the program immediately.", sig)
	cancel()
	_, ok = <-c
	if !ok {
		return
	}
	log.Println("Signal received twice, exit now.")
	os.Exit(0)
}

func newUnregisterFunc(c chan os.Signal, waitWatchSignal func()) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			signal.Stop(c)
			close(c)
		})
		waitWatchSignal()
	}
}
