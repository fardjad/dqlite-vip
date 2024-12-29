package utils

import (
	"sync"
	"time"
)

type WaitGroupWithTimeout struct {
	w sync.WaitGroup

	mu    sync.Mutex
	value int
}

func (w *WaitGroupWithTimeout) Add(delta int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.value += delta
	w.w.Add(delta)
}

func (w *WaitGroupWithTimeout) Done() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.value > 0 {
		w.value -= 1
		w.w.Done()
	}
}

func (w *WaitGroupWithTimeout) drain() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for w.value > 0 {
		w.value -= 1
		w.w.Done()
	}
}

// WaitWithTimeout waits for the WaitGroup to be done or until the duration has passed.
// Returns true if the WaitGroup is done, false otherwise.
func (w *WaitGroupWithTimeout) WaitWithTimeout(duration time.Duration) bool {
	done := make(chan struct{})

	go func() {
		w.w.Wait()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(duration):
		w.drain() // Forcefully drain if timeout occurs
		return false
	}
}
