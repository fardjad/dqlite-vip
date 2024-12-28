package time

import "time"

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type TickerFactoryFunc func(d time.Duration) Ticker
