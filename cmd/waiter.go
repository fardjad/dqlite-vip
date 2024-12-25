package cmd

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

type Waiter interface {
	Wait()
}

type SigTermWaiter struct{}

func (w *SigTermWaiter) Wait() {
	ch := make(chan os.Signal, 32)
	signal.Notify(ch, unix.SIGPWR)
	signal.Notify(ch, unix.SIGINT)
	signal.Notify(ch, unix.SIGQUIT)
	signal.Notify(ch, unix.SIGTERM)

	<-ch
}
