package utils

import "time"

type BetterTicker interface {
	C() <-chan time.Time
	Stop()
}

type BetterTickerFactoryFunc func(d time.Duration) BetterTicker

// Implements [BetterTicker]
type betterTicker struct {
	t *time.Ticker
}

func (t *betterTicker) C() <-chan time.Time {
	return t.t.C
}

func (t *betterTicker) Stop() {
	t.t.Stop()
}

// Implements [BetterTickerFactoryFunc]
func NewBetterTicker(d time.Duration) BetterTicker {
	return &betterTicker{t: time.NewTicker(d)}
}

// Implements [BetterTicker]
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
