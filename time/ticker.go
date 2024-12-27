package time

import "time"

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type realTicker struct {
	t *time.Ticker
}

func (t *realTicker) C() <-chan time.Time {
	return t.t.C
}

func (t *realTicker) Stop() {
	t.t.Stop()
}

type FakeTicker struct {
	c chan time.Time
}

func (t *FakeTicker) C() <-chan time.Time {
	return t.c
}

func (t *FakeTicker) Stop() {
}

func (t *FakeTicker) Tick(now time.Time) {
	t.c <- now
}

func NewFakeTicker() *FakeTicker {
	return &FakeTicker{c: make(chan time.Time)}
}

func NewRealTicker(d time.Duration) Ticker {
	return &realTicker{t: time.NewTicker(d)}
}
